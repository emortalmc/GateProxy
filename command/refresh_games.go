package command

import (
	"context"
	"go.minekube.com/brigodier"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"simple-proxy/luckperms"
	"simple-proxy/redisdb"
)

func newRefreshGamesCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {

	return brigodier.Literal("refreshgames").
		Requires(command.Requires(func(c *command.RequiresContext) bool {
			return luckperms.HasPermission(c.Source, "divine.refreshgames")
		})).
		Executes(command.Command(func(c *command.Context) error {
			redisdb.RedisClient.Publish(context.Background(), "proxyhello", "")

			return c.Source.SendMessage(&Text{
				Content: "Refreshed",
			})
		}))
}
