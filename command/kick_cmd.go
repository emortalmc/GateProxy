package command

import (
	"fmt"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"simple-proxy/luckperms"
	"simple-proxy/minimessage"
)

func newKickCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	var purple, _ = color.Make(color.LightPurple)
	var gold, _ = color.Make(color.Gold)
	var aqua, _ = color.Make(color.Aqua)

	return brigodier.Literal("kick").
		Requires(command.Requires(func(c *command.RequiresContext) bool {
			return luckperms.HasPermission(c.Source, "divine.kick")
		})).
		Executes(command.Command(func(c *command.Context) error {
			c.Source.SendMessage(&Text{
				Content: "Usage: /kick <player> <message>",
				S:       Style{Color: color.Gold},
			})
			return nil
		})).
		Then(
			brigodier.Argument("player", brigodier.String).Then(
				brigodier.Argument("message", brigodier.StringPhrase).
					Executes(command.Command(func(c *command.Context) error {
						playerStr := c.String("player")
						plr := p.PlayerByName(playerStr)
						if plr == nil {
							c.Source.SendMessage(&Text{
								S:       Style{Color: color.Red},
								Content: "Invalid player",
							})
							return nil
						}
						message := c.String("message")
						if message == "" {
							message = "get off my server"
						}

						plr.Disconnect(&Text{
							Extra: []Component{
								minimessage.Gradient("EmortalMC\n\n", Style{Bold: True}, *gold, *purple, *aqua),
								&Text{
									S:       Style{Color: color.Red},
									Content: fmt.Sprintf("You were kicked!\nReason: %s", message),
								},
							},
						})

						return nil
					})),
			),
		)
}
