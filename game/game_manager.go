package game

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"simple-proxy/redisdb"

	"go.minekube.com/common/minecraft/color"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/netutil"
	"go.minekube.com/gate/pkg/util/uuid"

	. "go.minekube.com/common/minecraft/component"
)

var (
	ctx     = context.Background()
	GameMap = make(map[string]string)
)

func RegisterPubSub(p *proxy.Proxy) {
	registerPubSub := redisdb.RedisClient.Subscribe(ctx, "registergame")
	ch := registerPubSub.Channel()

	go func() {
		for msg := range ch {
			fmt.Println(msg.Channel, msg.Payload)

			registerGame(p, msg.Payload)
		}
	}()

	joinPubSub := redisdb.RedisClient.Subscribe(ctx, "joingame")
	ch2 := joinPubSub.Channel()

	go func() {
		for msg := range ch2 {
			fmt.Println(msg.Channel, msg.Payload)

			args := strings.Fields(msg.Payload)
			game := args[0]
			pUuid := args[1]

			server := GameMap[game]
			if server == "" {
				continue
			}

			realID, err := uuid.Parse(pUuid)
			if err != nil {
				continue
			}

			player := p.Player(realID)
			if player == nil {
				continue
			}

			SendToServer(p, player, server, game, false, uuid.Nil)
		}
	}()

	result := redisdb.RedisClient.Publish(ctx, "proxyhello", "")
	if result.Err() != nil {
		log.Fatal("Redis failed to connect. Stopping...")
	}
}

func registerGame(p *proxy.Proxy, payload string) {
	// this is horrible, but I can't think of a better way :D
	localAddr := strings.Split(p.Servers()[0].ServerInfo().Addr().String(), ":")[0]

	args := strings.Fields(payload)
	gameName := strings.ToLower(args[0])
	serverName := strings.ToLower(args[1])
	port := args[2]

	addr, _ := netutil.Parse(localAddr+":"+port, "tcp")
	info := proxy.NewServerInfo(serverName, addr)

	log.Printf("Registered game:%s server:%s port:%s", gameName, serverName, port)

	GameMap[gameName] = serverName

	// TODO: reconnect players in limbo

	p.Register(info)
}

func SendToServer(p *proxy.Proxy, player proxy.Player, serverName string, game string, spectate bool, playerToSpectate uuid.UUID) {
	current := player.CurrentServer()
	if current == nil {
		log.Println("Not in a server")
		player.SendMessage(&Text{
			Extra: []Component{
				&Text{
					Content: "Not in a server",
					S:       Style{Color: color.Red, Bold: True},
				},
			},
		})
		return
	}

	if spectate {
		username := p.Player(playerToSpectate).Username()
		player.SendActionBar(&Text{
			Extra: []Component{
				&Text{
					Content: fmt.Sprintf("Spectating %s...", username),
					S:       Style{Color: color.Green},
				},
			},
		})
	} else {
		player.SendActionBar(&Text{
			Extra: []Component{
				&Text{
					Content: fmt.Sprintf("Joining %s...", game),
					S:       Style{Color: color.Green},
				},
			},
		})
	}

	//_ = player.WritePacket(&packet.EntitySoundEffect{
	//	SoundID:       34, // minecraft:block.amethyst_cluster.step
	//	SoundCategory: 0,
	//	EntityID:      packet.EntityStore.EntityID(player.ID()),
	//	Volume:        1,
	//	Pitch:         2,
	//	Seed:          0,
	//})

	if current.Server().ServerInfo().Name() == serverName {
		if spectate {
			go redisdb.RedisClient.Publish(ctx, "playerpubsub"+serverName, fmt.Sprintf("spectateplayer %s %s", player.ID(), playerToSpectate))
		} else {
			go redisdb.RedisClient.Publish(ctx, "playerpubsub"+serverName, fmt.Sprintf("changegame %s %s", player.ID(), game))
		}
		return
	}

	server := p.Server(serverName)
	if server == nil {
		log.Printf("Couldn't find %s", serverName)
		player.SendMessage(&Text{
			Extra: []Component{
				&Text{
					Content: fmt.Sprintf("Couldn't find %s", serverName),
					S:       Style{Color: color.Red, Bold: True},
				},
			},
		})
		return
	}

	var subgameStr string
	if spectate {
		subgameStr = fmt.Sprintf("%s %t %s", game, spectate, playerToSpectate)
	} else {
		subgameStr = game
	}
	res := redisdb.RedisClient.SetEX(ctx, fmt.Sprintf("%s-subgame", player.ID()), subgameStr, time.Second*10)
	if res.Err() != nil {
		log.Println("Failed to set subgame")
		player.SendMessage(&Text{
			Extra: []Component{
				&Text{
					Content: "Failed to join game",
					S:       Style{Color: color.Red},
				},
			},
		})
		return
	}

	_, err := player.CreateConnectionRequest(server).Connect(player.Context())
	if err != nil {
		log.Printf("Failed to join game. %s", err.Error())
		player.SendMessage(&Text{
			Extra: []Component{
				&Text{
					Content: fmt.Sprintf("Failed to join game! Error: %s", err.Error()),
					S:       Style{Color: color.Red},
				},
			},
		})
	}
}
