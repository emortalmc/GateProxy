package command

import (
	"fmt"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/uuid"
	"log"
	"simple-proxy/game"
	"strings"
)

func newSendallCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("sendall").
		Requires(command.Requires(func(c *command.RequiresContext) bool {
			return c.Source.HasPermission("divine.broadcast")
		})).
		Executes(command.Command(func(c *command.Context) error {
			c.Source.SendMessage(&Text{
				Content: "Usage: /sendall <game>",
				S:       Style{Color: color.Gold},
			})
			return nil
		})).
		Then(
			brigodier.Argument("game", brigodier.String).
				//Requires(command.Requires(func(c *command.RequiresContext) bool {
				//	return c.Source.HasPermission("divine.sendall")
				//})).
				Suggests(command.SuggestFunc(func(c *command.Context, b *brigodier.SuggestionsBuilder) *brigodier.Suggestions {
					for k := range game.GameMap {
						if strings.HasPrefix(k, b.RemainingLowerCase) {
							b.Suggest(k)
						}
					}
					return b.Build()
				})).
				Executes(command.Command(func(c *command.Context) error {
					player, ok := c.Source.(proxy.Player)

					if ok && player.Username() != "emortaldev" {
						c.Source.SendMessage(&Text{
							Content: "This command is restricted to epic gamers only",
						})
						return nil
					}

					senderName := "CONSOLE"
					if ok {
						senderName = player.Username()
					}

					gameName := c.String("game")
					server := game.GameMap[gameName]
					log.Printf("Server is %s for game %s", server, gameName)
					if server == "" {
						return nil
					}

					for _, plr := range p.Players() {
						go plr.SendMessage(&Text{
							Content: fmt.Sprintf("You were sent into %s by %s", gameName, senderName),
							S:       Style{Color: color.LightPurple},
						})
						go game.SendToServer(p, plr, server, gameName, false, uuid.Nil)
					}

					return nil
				})),
		)
}
