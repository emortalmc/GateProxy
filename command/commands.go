package command

import (
	"github.com/lucasb-eyer/go-colorful"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/common/minecraft/component/codec/legacy"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"math"
)

func RegisterCommands(p *proxy.Proxy) {
	p.Command().Register(newPlayCmd(p))
	registerWithAlias(p, newLobbyCmd(p), "l")
	p.Command().Register(newBroadcastCmd(p))
	p.Command().Register(newPingCmd(p))
	p.Command().Register(newDiscordCmd(p))
	p.Command().Register(newNbsCmd(p))
}

func registerWithAlias(p *proxy.Proxy, cmd brigodier.LiteralNodeBuilder, aliases ...string) {
	p.Command().Register(cmd)
	for _, alias := range aliases {
		p.Command().Register(brigodier.Literal(alias).Redirect(cmd.BuildLiteral()))
	}
}

var legacyCodec = &legacy.Legacy{Char: legacy.AmpersandChar}

func Gradient(content string, style Style, colors ...color.RGB) *Text {
	var component []Component
	chars := []rune(content)

	for i := range content {
		t := float64(i) / float64(len(content))

		hex, _ := color.Hex(LerpColor(t, colors...).Hex())

		style.Color = hex

		component = append(component, &Text{
			Content: string(chars[i]),
			S:       style,
		})
	}

	return &Text{
		Extra: component,
	}
}

func LerpColor(t float64, colors ...color.RGB) colorful.Color {
	t = math.Min(t, 1)

	if t == 1 {
		return colorful.Color(colors[len(colors)-1])
	}

	colorT := t * float64(len(colors)-1)
	newT := colorT - math.Floor(colorT)
	lastColor := colors[int(colorT)]
	nextColor := colors[int(colorT+1)]

	return colorful.Color{
		R: LerpInt(newT, nextColor.R, lastColor.R), G: LerpInt(newT, nextColor.G, lastColor.G), B: LerpInt(newT, nextColor.B, lastColor.B),
	}
}

func LerpInt(t float64, a float64, b float64) float64 {
	return a*t + b*(1-t)
}
