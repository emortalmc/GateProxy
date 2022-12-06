package command

import (
	"context"
	"fmt"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/uuid"
	"log"
	"math"
	"os"
	"simple-proxy/nbs"
	"strings"
	"time"
)

var magic_sound_events = [...]int{
	774,
	769,
	768,
	777,
	775,
	773,
	772,
	770,
	771,
	778,
	779,
	780,
	781,
	782,
	783,
	776,
	294,
	294,
}

var playingMap = make(map[uuid.UUID]context.CancelFunc)

func newNbsCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	files, err := os.ReadDir("./nbssongs/")
	if err != nil {
		log.Fatal(err)
	}

	return brigodier.Literal("nbs").
		Then(
			brigodier.Argument("songname", brigodier.StringPhrase).
				Suggests(command.SuggestFunc(func(c *command.Context, b *brigodier.SuggestionsBuilder) *brigodier.Suggestions {

					for _, file := range files {
						if strings.HasPrefix(file.Name(), b.Remaining) {
							b.Suggest(strings.Split(file.Name(), ".")[0])
						}
					}

					return b.Build()
				})).
				Executes(command.Command(func(c *command.Context) error {
					player, ok := c.Source.(proxy.Player)
					if !ok {
						return c.Source.SendMessage(&Text{Content: "Pong!"})
					}

					arg := c.String("songname")
					nbs, err := nbs.Read(fmt.Sprintf("./nbssongs/%s.nbs", arg))
					if err != nil {
						return err
					}

					prevCancel := playingMap[player.ID()]
					if prevCancel != nil {
						defer prevCancel()
					}

					i := 0
					ctx, cancel := context.WithCancel(player.Context())
					playingMap[player.ID()] = cancel
					go tickB(ctx, time.Millisecond*(time.Duration(1000/nbs.Tps)), func() {
						if i >= len(nbs.Ticks) {
							delete(playingMap, player.ID())
							defer cancel()
							return
						}

						tick := nbs.Ticks[i]

						green, _ := color.Make(color.Green)
						red, _ := color.Make(color.Red)
						outside, _ := color.Hex(LerpColor(float64(len(tick.Notes))/13, *green, *red).Hex())
						player.SendActionBar(Gradient(strings.Repeat("|", len(tick.Notes)*2), Style{}, *outside, *green, *outside))

						for _, note := range tick.Notes {
							_ = player.WritePacket(&EntitySoundEffect{
								SoundID:       magic_sound_events[int(math.Min(float64(note.Instrument), float64(len(magic_sound_events)-1)))],
								SoundCategory: 0,
								EntityID:      proxy.ServerConnectionEntityID(player.CurrentServer()),
								Volume:        float32(note.Volume) / 100,
								Pitch:         float32(math.Pow(2, (float64(note.Key)-float64(45))/float64(12))),
								Seed:          0,
							})
						}
						i++
					})

					return nil
				})),
		)
}

func tickB(ctx context.Context, interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for true {
		select {
		case <-ticker.C:
			fn()
		case <-ctx.Done():
			return
		}
	}
}
