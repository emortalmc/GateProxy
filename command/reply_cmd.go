package command

import (
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"simple-proxy/relationship"
)

func newReplyCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("reply").
		Executes(command.Command(func(c *command.Context) error {
			c.Source.SendMessage(&Text{
				Content: "Usage: /reply <message>",
				S:       Style{Color: color.Gold},
			})
			return nil
		})).
		Then(
			brigodier.Argument("message", brigodier.StringPhrase).
				Executes(command.Command(func(c *command.Context) error {
					player, ok := c.Source.(proxy.Player)
					if !ok {
						c.Source.SendMessage(&Text{
							Content: "Reply command cannot be used from console",
						})
						return nil
					}

					message := c.String("message")
					if message == "" {
						c.Source.SendMessage(&Text{
							S:       Style{Color: color.Red},
							Content: "Please type a message!",
						})
						return nil
					}

					lastMsg := relationship.LastMessageMap[player.ID()]
					plr := p.Player(lastMsg)
					if plr == nil {
						c.Source.SendMessage(&Text{
							S:       Style{Color: color.Red},
							Content: "That player went offline",
						})
						delete(relationship.LastMessageMap, player.ID())
					}

					SendMessage(message, player, plr)

					return nil
				})),
		)
}
