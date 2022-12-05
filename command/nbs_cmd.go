package command

import (
	"context"
	"fmt"
	"go.minekube.com/brigodier"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"math"
	"simple-proxy/nbs"
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

func newNbsCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("nbs").
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				return c.Source.SendMessage(&Text{Content: "Pong!"})
			}

			// TODO: Suggestions + argument

			nbs, err := nbs.Read("./Boss 2.nbs")
			if err != nil {
				return err
			}

			i := 0
			go tickB(player.Context(), int(nbs.Length), time.Millisecond*(time.Duration(1000/nbs.Tps)), func() {
				if i >= len(nbs.Ticks) {
					return
				}

				tick := nbs.Ticks[i]

				go player.SendActionBar(&Text{
					Content: fmt.Sprintf("tick: %d, notes: %d", i, len(tick.Notes)),
				})

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
		}))
}

func tickB(ctx context.Context, ticks int, interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	i := 0
	for i < ticks {
		select {
		case <-ticker.C:
			i++
			fn()
		case <-ctx.Done():
			return
		}
	}
}
