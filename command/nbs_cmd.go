package command

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"simple-proxy/minimessage"
	"simple-proxy/nbs"
	"simple-proxy/packet"

	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/util/uuid"
)

var magicSoundEvents = [...]int{ // all the note block sound event ids
	839, // minecraft:block.note_block.harp
	834, // minecraft:block.note_block.bass
	833, // minecraft:block.note_block.basedrum
	842, // minecraft:block.note_block.snare
	840, // minecraft:block.note_block.hat
	838, // minecraft:block.note_block.guitar
	837, // minecraft:block.note_block.flute
	835, // minecraft:block.note_block.bell
	836, // minecraft:block.note_block.chime
	843, // minecraft:block.note_block.xylophone
	844, // minecraft:block.note_block.iron_xylophone
	845, // minecraft:block.note_block.cow_bell
	846, // minecraft:block.note_block.didgeridoo
	847, // minecraft:block.note_block.bit
	848, // minecraft:block.note_block.banjo
	841, // minecraft:block.note_block.pling
}

var playingMap = make(map[uuid.UUID]context.CancelFunc)

func newNbsCmd(p *proxy.Proxy, tempAlias string) brigodier.LiteralNodeBuilder {
	creeperLyrics := map[int]string{
		32: "Creeper?",
		48: "Aw man",

		66:  "So we back in the mine",
		80:  "Got our pickaxe swinging from",
		96:  "Side",
		100: "to",
		104: "side",
		110: "Side-side",
		116: "to",
		120: "side",
		130: "This task, a grueling one",
		144: "Hope to find some diamonds",
		154: "tonight",
		164: "night", //
		174: "Diamonds tonight",

		192: "Heads up",
		206: "You hear a sound,",
		216: "turn around and",
		224: "look up",
		240: "Total shock fills your body",
		258: "Oh no, it's you again",
		272: "I can never forget those",
		288: "eyes",
		292: "eyes x2",
		296: "eyes x3",
		302: "eyes x4",
		306: "eyes x5",
		312: "'Cause, baby, tonight",

		332: "The creeper's trying to steal all our stuff again",
		374: "'Cause, baby, tonight",
		396: "Grab your pick, shovel, and bolt again",
		430: "Bolt again-gain",
		448: "And run, run until it's done, done",
		472: "Until the sun comes up in the 'morn",
		502: "'Cause, baby, tonight",
		524: "The creeper's trying to steal all our stuff again",
		558: "Stuff again-gain",
		578: "Just when you think you're safe",
		592: "Overhear some hissing from",
		608: "Right behind",
		620: "Right-right behind",
		642: "That's a nice life you have",
		656: "Shame it's gotta end at this",
		672: "time ",
		676: "time x2",
		680: "time x3",
		686: "time x4",
		690: "time x5",
		696: "time x6",

		704: "Blows up",
		718: "Then your health bar drops",
		728: "and you could use a",
		736: "one-up",
		752: "Get inside,", //
		760: "don't be tardy",
		770: "So now you're stuck in there",
		784: "Half a heart is left,",
		794: "but don't",
		800: "die",
		804: "die x2",
		808: "die x3",
		814: "die x4",
		818: "die x5",
		824: "'Cause, baby, tonight",

		844:  "The creeper's trying to steal all our stuff again",
		886:  "'Cause, baby, tonight",
		908:  "Grab your pick, shovel, and bolt again",
		942:  "Bolt again-gain",
		960:  "And run, run until it's done, done",
		984:  "Until the sun comes up in the 'morn",
		1014: "'Cause, baby, tonight",
		1036: "The creeper's trying to steal all our stuff again",
		1068: "Stuff again-gain",

		// I can't be arsed to caption the rapping part
		1088: "(I don't wanna caption this)",

		1286: "'Cause, baby, tonight",
		1308: "The creeper's trying to steal all our stuff again",
		1350: "Yeah, baby, tonight",
		1374: "Grab your sword, armour and go",
		1405: "Take your revenge",
		1420: "So fight, fight",
		1432: "like it's the",
		1440: "last, last night of your",
		1456: "life, life",
		1464: "Show them your bite",
		1476: "'Cause, baby, tonight",
		1500: "The creeper's trying to steal all our stuff again",
		1542: "'Cause, baby, tonight",
		1564: "Grab your pick, shovel, and bolt again",
		1598: "Bolt again-gain",
		1616: "And run, run until it's done, done",
		1640: "Until the sun comes up in the 'morn",
		1670: "'Cause, baby, tonight",
		1692: "The creeper's trying to steal all our stuff again",
	}

	files, err := os.ReadDir("./nbssongs/")
	if err != nil {
		fmt.Println("nbssongs folder doesn't exist")
	}

	return brigodier.Literal(tempAlias).
		Then(
			brigodier.Argument("songname", brigodier.StringPhrase).
				Suggests(command.SuggestFunc(func(c *command.Context, b *brigodier.SuggestionsBuilder) *brigodier.Suggestions {

					for _, file := range files {
						if strings.HasPrefix(file.Name(), b.Remaining) {
							b.Suggest(strings.Split(file.Name(), ".")[0])
						}
					}

					b.Suggest("stop")

					return b.Build()
				})).
				Executes(command.Command(func(c *command.Context) error {
					player, ok := c.Source.(proxy.Player)
					if !ok {
						return c.Source.SendMessage(&Text{Content: "Pong!"})
					}

					arg := c.String("songname")

					prevCancel := playingMap[player.ID()]
					if prevCancel != nil {
						defer prevCancel()
					}

					nbs, err := nbs.Read(fmt.Sprintf("./nbssongs/%s.nbs", arg))
					if err != nil {
						return nil
					}

					second, _ := color.Hex("#e64ce6")
					third, _ := color.Hex("#009dff")
					authorStr := nbs.OriginalAuthor
					if authorStr == "" {
						authorStr = nbs.Author
					}
					titleGradient := minimessage.Gradient(fmt.Sprintf("%s - %s", authorStr, nbs.Name), Style{}, *second, *third)
					player.SendActionBar(titleGradient)

					i := 0
					ctx, cancel := context.WithCancel(player.Context())
					playingMap[player.ID()] = cancel

					tickDuration := time.Millisecond * (time.Duration(1000.0 / nbs.Tps))
					songDuration := int64((tickDuration * time.Duration(nbs.Length)) / time.Second)
					songStart := time.Now()

					go tickB(ctx, tickDuration, func() {
						if i >= len(nbs.Ticks) {
							delete(playingMap, player.ID())
							defer cancel()
							return
						}

						tick := nbs.Ticks[i]

						if nbs.Name == "DJ Got Us Fallin' in Love" {
							lyric := creeperLyrics[i]
							if lyric != "" {
								player.SendActionBar(&Text{Content: lyric})
							}
						} else {
							if uint16(i)%nbs.Tps == 0 {
								elapsed := int64(time.Since(songStart) / time.Second)

								player.SendActionBar(&Text{Extra: []Component{
									titleGradient,
									&Text{
										S:       Style{Color: color.Gray},
										Content: " | ",
									},
									&Text{
										Content: fmt.Sprintf("%s / %s", formatTime(elapsed), formatTime(songDuration)),
									},
								}})
							}
						}
						//green, _ := minimessage.Make(minimessage.Green)
						//red, _ := minimessage.Make(minimessage.Red)
						//outside, _ := minimessage.Hex(LerpColor(float64(len(tick.Notes))/13, *green, *red).Hex())
						//player.SendActionBar(Gradient(strings.Repeat("|", len(tick.Notes)*2), Style{}, *outside, *green, *outside))

						for _, note := range tick.Notes {
							_ = player.WritePacket(&packet.EntitySoundEffect{
								SoundID:       magicSoundEvents[int(math.Min(float64(note.Instrument), float64(len(magicSoundEvents)-1)))],
								SoundCategory: 0,
								EntityID:      packet.EntityStore.EntityID(player.ID()),
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

func formatTime(secs int64) string {
	minutes := secs / 60
	seconds := secs % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
