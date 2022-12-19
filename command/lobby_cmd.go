package command

import (
	"go.minekube.com/brigodier"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/uuid"
	"simple-proxy/game"
)

func newLobbyCmd(p *proxy.Proxy, tempAlias string) brigodier.LiteralNodeBuilder {
	return brigodier.Literal(tempAlias).
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				go c.Source.SendMessage(&Text{
					Content: "Lobby command cannot be used from console",
				})
				return nil
			}

			game.SendToServer(p, player, "lobby", "lobby", false, uuid.Nil)

			return nil
		}))
}
