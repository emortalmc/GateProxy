package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"simple-proxy/command"
	"simple-proxy/game"
	"simple-proxy/luckperms"
	"simple-proxy/minimessage"
	"simple-proxy/packet"
	"simple-proxy/redisdb"
	"simple-proxy/webhook"

	"go.minekube.com/gate/pkg/edition/java/ping"
	"go.minekube.com/gate/pkg/edition/java/proto/state"
	"go.minekube.com/gate/pkg/edition/java/proto/version"
	"go.minekube.com/gate/pkg/edition/java/proxy/tablist"

	"github.com/robinbraemer/event"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var discordWebhookURL string

var purple, _ = color.Make(color.LightPurple)
var gold, _ = color.Make(color.Gold)
var aqua, _ = color.Make(color.Aqua)
var ip, _ = color.Hex("#266ee0")

func main() {
	const discordEnv = "DISCORD_WEBHOOK_URL"
	discordWebhookURL = os.Getenv(discordEnv)
	if discordWebhookURL == "" {
		_, _ = fmt.Fprintln(os.Stderr, discordEnv)
		os.Exit(1)
	}

	redisdb.RedisClient = redisdb.InitRedis()

	state.Play.ClientBound.Register(&packet.EntitySoundEffect{}, &state.PacketMapping{
		ID:       0x5F,
		Protocol: version.Minecraft_1_19_1.Protocol,
	})

	proxy.Plugins = append(proxy.Plugins, proxy.Plugin{
		Name: "SimpleProxy",
		Init: func(ctx context.Context, proxy *proxy.Proxy) error {
			return newSimpleProxy(proxy).init()
		},
	})

	// Execute Gate entrypoint and block until shutdown.
	// We could also run gate.Start if we don't need Gate's command-line.
	gate.Execute()
}

type SimpleProxy struct {
	*proxy.Proxy
}

func newSimpleProxy(proxy *proxy.Proxy) *SimpleProxy {
	return &SimpleProxy{
		Proxy: proxy,
	}
}

func (p *SimpleProxy) init() error {
	command.RegisterCommands(p.Proxy)
	p.registerSubscribers()

	packet.EntityStore.Subscribe(p.Event())
	game.RegisterPubSub(p.Proxy)
	return nil
}

// Register event subscribers
func (p *SimpleProxy) registerSubscribers() {
	// Send message on server switch.
	event.Subscribe(p.Event(), 0, p.onServerLogin)
	event.Subscribe(p.Event(), 0, p.onServerDisconnect)
	event.Subscribe(p.Event(), 0, p.onChat)

	// Change the MOTD response.
	event.Subscribe(p.Event(), 0, pingHandler(p.Proxy))
}

var randomJoinMessages = []string{
	"How's the wife?",
	"How's the kids?",
	"What's the weather like?",
	"Back so soon?",
	"A good day for EmortalMC",
	"Another great day for procrastination",
	"Great to see you!",
	"[Server] Back in 5 minutes",
	"Salutations, <username>",
	"Act busy, <username> is here",
	"Not you again...",
	"I hope you brought pizza",
	"I hope you brought friends",
	"I hope you aren't using Lunar",
	"Welcome back, we missed you",
	"You finally arrived!",
	"Marathon again?",
}

func (p *SimpleProxy) onServerLogin(e *proxy.ServerPostConnectEvent) {
	if e.PreviousServer() == nil {
		refreshTablist(p.Proxy)
		webhook.PlayerJoined(e.Player(), p.PlayerCount(), discordWebhookURL)
		collectResult := luckperms.CollectData(e.Player())
		if collectResult != nil {
			e.Player().SendMessage(&Text{
				Content: "Failed to collect your LuckPerms data!",
			})
			fmt.Printf("failed to collect %s's LuckPerms data %s\n", e.Player().Username(), collectResult)
		}

		thereAre := "There are now"
		plrs := "players"
		if p.PlayerCount() == 1 {
			plrs = "player"
			thereAre = "There is now"
		}
		e.Player().SendMessage(&Text{
			Extra: []Component{
				&Text{
					S:       Style{Color: color.Gray},
					Content: "Welcome to ",
				},
				minimessage.Gradient("EmortalMC", Style{Bold: True}, *gold, *purple),
				&Text{
					S:       Style{Color: color.Gray},
					Content: fmt.Sprintf("! %s ", thereAre),
				},
				&Text{
					S:       Style{Color: color.Yellow},
					Content: strconv.Itoa(p.PlayerCount()),
				},
				&Text{
					S:       Style{Color: color.Gray},
					Content: fmt.Sprintf(" %s online", plrs),
				},
			},
		})

		randomMessage := minimessage.Gradient(strings.Replace(randomJoinMessages[rand.Intn(len(randomJoinMessages))], "<username>", e.Player().Username(), 1), Style{}, *gold, *purple)

		ctx, cancel := context.WithCancel(e.Player().Context())
		i := 0
		go tick(ctx, 2*time.Second, func() {
			if i > 4 {
				defer cancel()
				return
			}
			i += 1

			e.Player().SendActionBar(randomMessage)
		})

	}

}

func (p *SimpleProxy) onServerDisconnect(e *proxy.DisconnectEvent) {
	refreshTablist(p.Proxy)
	webhook.PlayerLeft(e.Player(), p.PlayerCount(), discordWebhookURL)
}

func (p *SimpleProxy) onPreJoin(e *proxy.LoginEvent) {
	for _, a := range command.BanMap {
		if a == e.Player().ID() {
			e.Deny(&Text{
				Content: "get band",
			})
			break
		}
	}
}

func (p *SimpleProxy) onChat(e *proxy.PlayerChatEvent) {
	e.SetAllowed(false)

	components := []Component{
		luckperms.DisplayName(e.Player()),
		&Text{
			S:       Style{Color: color.Gray},
			Content: ": ",
		},
		&Text{
			S:       Style{Color: color.White},
			Content: e.Message(),
		},
	}

	for _, plr := range p.Players() {
		go plr.SendMessage(&Text{
			Extra: components,
		})
	}
}

func refreshTablist(p *proxy.Proxy) {

	for _, plr := range p.Players() {
		go tablist.SendHeaderFooter(plr,
			&Text{
				Extra: []Component{
					&Text{
						S:       Style{Color: gold},
						Content: "┌                                                  ",
					},
					&Text{
						S:       Style{Color: purple},
						Content: "┐ \n",
					},
					minimessage.Gradient("EmortalMC\n", Style{Bold: True}, *gold, *purple),
				},
			},
			&Text{
				Extra: []Component{
					&Text{
						S:       Style{Color: color.Gray},
						Content: fmt.Sprintf(" \n%d online", p.PlayerCount()),
					},
					&Text{
						S:       Style{Color: ip},
						Content: "\nᴍᴄ.ᴇᴍᴏʀᴛᴀʟ.ᴅᴇᴠ",
					},
					&Text{
						S:       Style{Color: purple},
						Content: "\n└                                                  ",
					},
					&Text{
						S:       Style{Color: gold},
						Content: "┘ ",
					},
				},
			},
		)
	}
}

func pingHandler(p *proxy.Proxy) func(evt *proxy.PingEvent) {
	messages := []string{
		"coolest server to ever exist",
		"better than hypixel",
		"you should join",
		"stop scrolling, click here!",
		"Lunar client users: Beware!",
		"using 3 server softwares!",
		"gradient lover",
		"emortal is watching",
		"Chuck Norris joined and said it was pretty good",
		"private lobbies when?",
	}

	return func(e *proxy.PingEvent) {
		randomMessage := messages[rand.Intn(len(messages))]

		motd := &Text{
			Extra: []Component{
				&Text{
					Content: "▓▒░              ",
					S:       Style{Color: color.LightPurple},
				},
				&Text{
					Content: "⚡   ",
					S:       Style{Color: color.LightPurple, Bold: True},
				},
				minimessage.Gradient("EmortalMC", Style{Bold: True}, *gold, *purple),
				&Text{
					Content: "   ⚡",
					S:       Style{Color: color.Gold, Bold: True},
				},
				&Text{
					Content: "              ░▒▓",
					S:       Style{Color: color.Gold},
				},
				&Text{
					Content: "\n" + randomMessage,
					S:       Style{Color: color.Yellow},
				},
			},
		}

		evt := e.Ping()
		evt.Description = motd
		evt.Players.Max = evt.Players.Online + 1

		sampleCount := int(math.Min(float64(p.PlayerCount()), 10))
		samples := make([]ping.SamplePlayer, sampleCount)
		for i, plr := range p.Players() {
			if i >= sampleCount {
				break
			}
			samples[i] = ping.SamplePlayer{Name: plr.Username(), ID: plr.ID()}
		}

		evt.Players.Sample = samples
	}
}

// tick runs a function every interval until the context is cancelled.
func tick(ctx context.Context, interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for true {
		select {
		case <-ticker.C:
			fn()
		case <-ctx.Done():
			return
		}
	}
}
