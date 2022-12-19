package command

import (
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"simple-proxy/luckperms"
	"simple-proxy/relationship"
)

var dmGreen, _ = color.Hex("#00BD5D")
var dmDarkGreen, _ = color.Hex("#2F8F49")
var dmTextColor, _ = color.Hex("#A5F0BA")

var youText = &Text{Content: "YOU", S: Style{Bold: True, Color: dmGreen}}
var arrowText = &Text{Content: " → ", S: Style{Color: color.Gray}}

func newMsgCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("msg").
		Executes(command.Command(func(c *command.Context) error {
			c.Source.SendMessage(&Text{
				Content: "Usage: /msg <player> <message>",
				S:       Style{Color: color.Gold},
			})
			return nil
		})).
		Then(
			brigodier.Argument("player", brigodier.String).Then(
				brigodier.Argument("message", brigodier.StringPhrase).
					Executes(command.Command(func(c *command.Context) error {
						player, ok := c.Source.(proxy.Player)
						if !ok {
							c.Source.SendMessage(&Text{
								Content: "Msg command cannot be used from console",
							})
							return nil
						}

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
							c.Source.SendMessage(&Text{
								S:       Style{Color: color.Red},
								Content: "Please type a message!",
							})
							return nil
						}

						SendMessage(message, player, plr)

						return nil
					})),
			),
		)
}

func SendMessage(message string, player proxy.Player, other proxy.Player) {
	msgText := &Text{Content: message, S: Style{Color: dmTextColor}}

	relationship.LastMessageMap[player.ID()] = other.ID()
	relationship.LastMessageMap[other.ID()] = player.ID()

	player.SendMessage(&Text{
		Extra: []Component{
			&Text{Content: "[", S: Style{Color: dmDarkGreen}},
			youText,
			arrowText,
			luckperms.DisplayName(other),
			&Text{Content: "] ", S: Style{Color: dmDarkGreen}},
			msgText,
		},
	})
	other.SendMessage(&Text{
		Extra: []Component{
			&Text{Content: "[", S: Style{Color: dmDarkGreen}},
			luckperms.DisplayName(player),
			arrowText,
			youText,
			&Text{Content: "] ", S: Style{Color: dmDarkGreen}},
			msgText,
		},
	})
}
