package command

import (
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

func newPlayCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("play").
		Executes(command.Command(func(c *command.Context) error {
			c.Source.SendMessage(&Text{
				Content: "Usage: /play <game>",
				S:       Style{Color: color.Gold},
			})
			return nil
		})).
		Then(
			brigodier.Argument("game", brigodier.String).
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

					game.SendToServer(p, player, server, gameName, false, uuid.Nil)

					return nil
				})),
		)
}
