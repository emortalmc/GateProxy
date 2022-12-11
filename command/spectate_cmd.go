package command

import (
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"simple-proxy/game"
	"strings"
)

func newSpectateCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("spectate").
		Executes(command.Command(func(c *command.Context) error {
			c.Source.SendMessage(&Text{
				Content: "Usage: /spectate <username>",
				S:       Style{Color: color.Gold},
			})
			return nil
		})).
		Then(
			brigodier.Argument("username", brigodier.String).
				Suggests(command.SuggestFunc(func(c *command.Context, b *brigodier.SuggestionsBuilder) *brigodier.Suggestions {
					for k := range game.GameMap {
						str := game.GameMap[k]
						if strings.HasPrefix(str, b.RemainingLowerCase) {
							b.Suggest(str)
						}
					}
					return b.Build()
				})).
				Executes(command.Command(func(c *command.Context) error {
					player, ok := c.Source.(proxy.Player)
					if !ok {
						c.Source.SendMessage(&Text{
							Content: "Spectate command cannot be used from console",
						})
						return nil
					}

					spectateName := c.String("username")
					plr := p.PlayerByName(spectateName)
					if plr == nil {
						c.Source.SendMessage(&Text{
							S:       Style{Color: color.Red},
							Content: "Could not find that player",
						})
						return nil
					}

					serverName := plr.CurrentServer().Server().ServerInfo().Name()

					game.SendToServer(p, player, serverName, "spectate", true, plr.ID())

					return nil
				})),
		)
}
