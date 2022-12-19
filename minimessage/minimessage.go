package minimessage

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/common/minecraft/component/codec/legacy"
	"math"
	"strings"
)

type tree struct {
	Parent   *tree
	Children []*tree
	Key      string
	Value    string
}

func Parse(miniString string) *Text {
	var styles []Style
	styles = append(styles, Style{Color: color.White})

	var components []Component

	for _, s := range strings.Split(miniString, "<") {
		if s == "" {
			continue
		}

		split := strings.Split(s, ">")

		key := split[0]
		if strings.HasPrefix(key, "/") {
			styles = styles[:len(styles)-1]
		} else {
			newStyle := styles[len(styles)-1]

			styles = append(styles, newStyle)
		}

		newText := Modify(key, split[1], &styles[len(styles)-1])
		components = append(components, newText)

	}

	return &Text{
		Extra: components,
	}
}

func Modify(key string, content string, style *Style) *Text {
	newText := &Text{}

	switch {
	case strings.HasPrefix(key, "#"): // <#ff00ff>
		parsed, err := ParseColor(key)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		style.Color = parsed
		newText.Content = content
		newText.S = *style
	case strings.HasPrefix(key, "color"): // <color:light_purple>
		colorName := strings.Split(key, ":")[1]
		parsed, err := ParseColor(colorName)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		style.Color = parsed
		newText.Content = content
		newText.S = *style

	case key == "bold": // <bold>
		style.Bold = True
		newText.Content = content
		newText.S = *style

	case strings.HasPrefix(key, "gradient"): // <gradient:light_purple:gold>
		colorKey := strings.Split(key, ":")
		colorNames := colorKey[1:]

		colors := make([]color.RGB, len(colorNames))
		for i, col := range colorNames {
			parsedColor, err := ParseColor(col)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			newColor, _ := color.Make(parsedColor)
			colors[i] = *newColor
		}

		newText = Gradient(content, *style, colors...)
	}

	return newText
}

func ParseColor(name string) (color.Color, error) {
	if strings.HasPrefix(name, "#") {
		newColor, err := color.Hex(name)
		if err != nil {
			return nil, err
		}
		return newColor, nil
	} else {
		return FromName(name)
	}
}

func FromName(name string) (color.Color, error) {
	for _, a := range color.Names {
		if a.String() == name {
			return a, nil
		}
	}
	return nil, fmt.Errorf("invalid minimessage name %s", name)
}

var LegacyCodec = &legacy.Legacy{Char: legacy.AmpersandChar}

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
