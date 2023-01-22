package command

import (
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func RegisterCommands(p *proxy.Proxy) {
	p.Command().Register(newPlayCmd(p))

	p.Command().Register(newLobbyCmd(p, "lobby"))
	p.Command().Register(newLobbyCmd(p, "l"))
	p.Command().Register(newLobbyCmd(p, "hub"))

	p.Command().Register(newBroadcastCmd(p))
	p.Command().Register(newPingCmd(p))
	p.Command().Register(newDiscordCmd(p))

	p.Command().Register(newNbsCmd(p, "nbs"))
	p.Command().Register(newNbsCmd(p, "music"))

	p.Command().Register(newMsgCmd(p))
	p.Command().Register(newReplyCmd(p))
	p.Command().Register(newSendallCmd(p))
	p.Command().Register(newSpectateCmd(p))
	p.Command().Register(newShrugCmd(p))

	p.Command().Register(newKickCmd(p))
	p.Command().Register(newRefreshGamesCmd(p))
}
