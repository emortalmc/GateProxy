package command

import (
	"fmt"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func newPingCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("ping").
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				return c.Source.SendMessage(&Text{Content: "Pong!"})
			}

			return player.SendMessage(&Text{
				Content: fmt.Sprintf("Pong! Your ping is %s", player.Ping()),
				S:       Style{Color: color.Green},
			})
		}))
}
