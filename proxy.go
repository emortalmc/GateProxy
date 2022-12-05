package main

import (
	"context"
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"go.minekube.com/gate/pkg/edition/java/ping"
	"go.minekube.com/gate/pkg/util/uuid"
	"log"
	"math"
	"math/rand"
	"simple-proxy/game"
	"time"

	"github.com/robinbraemer/event"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/common/minecraft/component/codec/legacy"
	"go.minekube.com/gate/cmd/gate"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func main() {
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

// SimpleProxy is a simple proxy to showcase some features of Gate.
//
// In this example:
//   - Add a `/broadcast` command
//   - Send a message when player switches the server
//   - Show boss bars to players
type SimpleProxy struct {
	*proxy.Proxy
}

var legacyCodec = &legacy.Legacy{Char: legacy.AmpersandChar}

func newSimpleProxy(proxy *proxy.Proxy) *SimpleProxy {
	return &SimpleProxy{
		Proxy: proxy,
	}
}

// initialize our sample proxy
func (p *SimpleProxy) init() error {
	p.registerCommands()
	p.registerSubscribers()

	game.RegisterPubSub(p.Proxy)
	return nil
}

// Register a proxy-wide commands (can be run while being on any server)
func (p *SimpleProxy) registerCommands() {

	p.Command().Register(brigodier.Literal("play").Executes(command.Command(func(c *command.Context) error {
		c.Source.SendMessage(&Text{
			Content: "Usage: /play <game>",
			S:       Style{Color: color.Gold},
		})
		return nil
	})).Then(
		brigodier.Argument("game", brigodier.String).
			Suggests(command.SuggestFunc(func(
				c *command.Context,
				b *brigodier.SuggestionsBuilder,
			) *brigodier.Suggestions {
				for k := range game.GameMap {
					b.Suggest(game.GameMap[k])
				}
				return b.Build()
			})).
			Executes(command.Command(func(c *command.Context) error {
				player, ok := c.Source.(proxy.Player)
				if !ok {
					go c.Source.SendMessage(&Text{
						Content: "Play command cannot be used from console",
					})
					return nil
				}

				gameName := c.String("game")
				server := game.GameMap[gameName]
				log.Printf("Server is %s for game %s", server, gameName)
				if server == "" {
					return nil
				}

				game.SendToServer(p.Proxy, player, server, gameName, false, uuid.Nil)

				return nil
			})),
	))

	// Registers the "/broadcast" command
	p.Command().Register(brigodier.Literal("broadcast").Then(
		// Adds message argument as in "/broadcast <message>"
		brigodier.Argument("message", brigodier.StringPhrase).
			// Adds completion suggestions as in "/broadcast [suggestions]"
			Suggests(command.SuggestFunc(func(
				c *command.Context,
				b *brigodier.SuggestionsBuilder,
			) *brigodier.Suggestions {
				player, ok := c.Source.(proxy.Player)
				if ok {
					b.Suggest("&oI am &6&l" + player.Username())
				}
				b.Suggest("Hello world!")
				return b.Build()
			})).
			// Executed when running "/broadcast <message>"
			Executes(command.Command(func(c *command.Context) error {
				// Colorize/format message
				message, err := legacyCodec.Unmarshal([]byte(c.String("message")))
				if err != nil {
					return c.Source.SendMessage(&Text{
						Content: fmt.Sprintf("Error formatting message: %v", err)})
				}

				// Send to all players on this proxy
				for _, player := range p.Players() {
					// Send message in new goroutine to not block
					// this loop if any player has a slow connection.
					go func(p proxy.Player) { _ = p.SendMessage(message) }(player)
				}
				return nil
			})),
	))
	p.Command().Register(brigodier.Literal("ping").
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				return c.Source.SendMessage(&Text{Content: "Pong!"})
			}
			return player.SendMessage(&Text{
				Content: fmt.Sprintf("Pong! Your ping is %s", player.Ping()),
				S:       Style{Color: color.Green},
			})
		})),
	)
}

// Register event subscribers
func (p *SimpleProxy) registerSubscribers() {
	// Send message on server switch.
	event.Subscribe(p.Event(), 0, p.onServerSwitch)

	// Change the MOTD response.
	event.Subscribe(p.Event(), 0, pingHandler(p.Proxy))
}

func (p *SimpleProxy) onServerSwitch(e *proxy.ServerPostConnectEvent) {
	newServer := e.Player().CurrentServer()
	if newServer == nil {
		return
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
	}

	return func(e *proxy.PingEvent) {
		randomMessage := messages[rand.Intn(len(messages))]

		first, _ := color.Make(color.Gold)
		second, _ := color.Make(color.LightPurple)

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
				gradient("EmortalMC", *first, *second),
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
	for {
		select {
		case <-ticker.C:
			fn()
		case <-ctx.Done():
			return
		}
	}
}

func gradient(content string, first color.RGB, second color.RGB) *Text {
	var component []Component
	chars := []rune(content)

	for i := range content {
		t := float64(i) / float64(len(content))

		hex, _ := color.Hex(lerpColor(t, first, second).Hex())

		component = append(component, &Text{
			Content: string(chars[i]),
			S:       Style{Color: hex, Bold: True},
		})
	}

	return &Text{
		Extra: component,
	}
}

func lerpColor(t float64, a color.RGB, b color.RGB) colorful.Color {
	return colorful.Color{
		R: lerpInt(t, a.R, b.R), G: lerpInt(t, a.G, b.G), B: lerpInt(t, a.B, b.B),
	}
}

func lerpInt(t float64, a float64, b float64) float64 {
	a1 := float64(a)
	b1 := float64(b)
	return a1*t + b1*(1-t)
}
