package command

import (
	"fmt"
	"go.minekube.com/brigodier"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func newBroadcastCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("broadcast").Then(
		brigodier.Argument("message", brigodier.StringPhrase).
			Executes(command.Command(func(c *command.Context) error {
				// Colorize/format message
				message, err := legacyCodec.Unmarshal([]byte(c.String("message")))
				if err != nil {
					return c.Source.SendMessage(&Text{
						Content: fmt.Sprintf("Error formatting message: %v", err)})
				}

				// Send to all players on this proxy
				for _, player := range p.Players() {
					// Send message in new goroutine to not block
					// this loop if any player has a slow connection.
					go player.SendMessage(message)
				}
				return nil
			})),
	)
}
