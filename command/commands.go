package command

import (
	"go.minekube.com/brigodier"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func RegisterCommands(p *proxy.Proxy) {
	p.Command().Register(newPlayCmd(p))
	registerWithAlias(p, newLobbyCmd(p), "l")
	p.Command().Register(newBroadcastCmd(p))
	p.Command().Register(newPingCmd(p))
	p.Command().Register(newDiscordCmd(p))
	p.Command().Register(newNbsCmd(p))
	p.Command().Register(newMsgCmd(p))
	p.Command().Register(newReplyCmd(p))
	p.Command().Register(newSendallCmd(p))
	p.Command().Register(newSpectateCmd(p))
	p.Command().Register(newShrugCmd(p))
}

func registerWithAlias(p *proxy.Proxy, cmd brigodier.LiteralNodeBuilder, aliases ...string) {
	p.Command().Register(cmd)
	for _, alias := range aliases {
		p.Command().Register(brigodier.Literal(alias).Redirect(cmd.BuildLiteral()))
	}
}
