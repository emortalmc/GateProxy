package command

import (
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func newDiscordCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	color1, _ := color.Hex("#7289da")
	color2, _ := color.Hex("#51629c")

	return brigodier.Literal("discord").
		Executes(command.Command(func(c *command.Context) error {

			return c.Source.SendMessage(&Text{
				S: Style{
					HoverEvent: ShowText(&Text{Content: "https://discord.gg/TZyuMSha96", S: Style{Color: color.Gray}}),
					ClickEvent: OpenUrl("https://discord.gg/TZyuMSha96"),
				},
				Extra: []Component{
					Gradient("Click to join our", Style{}, *color1, *color2, *color2),
					&Text{
						S:       Style{Bold: True, Color: color1},
						Content: " Discord",
					},
					&Text{
						S:       Style{Color: color2},
						Content: "!",
					},
				},
			})
		}))
}
