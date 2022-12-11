package nbs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type NBS struct {
	version            byte
	instruments        byte
	Length             uint16
	layers             uint16
	Name               string
	Author             string
	OriginalAuthor     string
	description        string
	Tps                uint16
	autoSaving         bool
	autoSavingDuration byte
	timeSignature      byte
	minutesSpent       int32
	leftClicks         int32
	rightClicks        int32
	notesAdded         int32
	notesRemoved       int32
	midiName           string
	loop               bool
	maxLoops           byte
	loopStart          uint16

	Ticks []Tick
}

type Tick struct {
	Tick  uint16
	Notes []Note
}

type Note struct {
	Instrument byte
	Key        byte
	Volume     byte
	Pan        byte
	Pitch      uint16
}

func Read(path string) (NBS, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return NBS{}, err
	}

	buf := bytes.NewBuffer(file)

	// First two bytes are always zero
	buf.ReadByte()
	buf.ReadByte()

	nbs := NBS{}

	version, _ := buf.ReadByte()
	instruments, _ := buf.ReadByte()
	var length uint16
	_ = binary.Read(buf, binary.LittleEndian, &length)
	var layers uint16
	_ = binary.Read(buf, binary.LittleEndian, &layers)

	name, _ := readString(buf)
	author, _ := readString(buf)
	originalAuthor, _ := readString(buf)
	description, _ := readString(buf)

	var tps uint16
	_ = binary.Read(buf, binary.LittleEndian, &tps)

	autoSaving, _ := buf.ReadByte()
	autoSavingDuration, _ := buf.ReadByte()
	timeSignature, _ := buf.ReadByte()

	var minutesSpent int32
	_ = binary.Read(buf, binary.LittleEndian, &minutesSpent)
	var leftClicks int32
	_ = binary.Read(buf, binary.LittleEndian, &leftClicks)
	var rightClicks int32
	_ = binary.Read(buf, binary.LittleEndian, &rightClicks)
	var notesAdded int32
	_ = binary.Read(buf, binary.LittleEndian, &notesAdded)
	var notesRemoved int32
	_ = binary.Read(buf, binary.LittleEndian, &notesRemoved)

	midiName, _ := readString(buf)
	loop, _ := buf.ReadByte()
	maxLoops, _ := buf.ReadByte()

	var loopStart uint16
	_ = binary.Read(buf, binary.LittleEndian, &loopStart)

	nbs.version = version
	nbs.instruments = instruments
	nbs.Length = length
	nbs.layers = layers
	nbs.Name = name
	nbs.Author = author
	nbs.OriginalAuthor = originalAuthor
	nbs.description = description
	nbs.Tps = tps / 100
	nbs.autoSaving = autoSaving == 1
	nbs.autoSavingDuration = autoSavingDuration
	nbs.timeSignature = timeSignature
	nbs.minutesSpent = minutesSpent
	nbs.leftClicks = leftClicks
	nbs.rightClicks = rightClicks
	nbs.notesAdded = notesAdded
	nbs.notesRemoved = notesRemoved
	nbs.midiName = midiName
	nbs.loop = loop == 1
	nbs.maxLoops = maxLoops
	nbs.loopStart = loopStart

	nbs.Ticks = readNotes(buf, nbs)

	fmt.Printf("\n\nRead %d ticks \n\n", len(nbs.Ticks))

	return nbs, nil
}

func readString(buf *bytes.Buffer) (string, error) {
	var lengthOfUtf int32
	_ = binary.Read(buf, binary.LittleEndian, &lengthOfUtf)
	utfBytes := make([]byte, lengthOfUtf)
	_, err := buf.Read(utfBytes)
	if err != nil {
		return "", err
	}

	return string(utfBytes), nil
}

func readNotes(buf *bytes.Buffer, nbs NBS) []Tick {
	tick := -1
	ticks := make([]Tick, nbs.Length+1)

	for true {
		var jumps uint16
		_ = binary.Read(buf, binary.LittleEndian, &jumps)

		if jumps == 0 {
			break
		}

		tick += int(jumps)

		notes := readTickNoteLayers(buf)

		ticks[tick] = Tick{
			Tick:  uint16(tick),
			Notes: notes,
		}
	}

	return ticks
}

func readTickNoteLayers(buf *bytes.Buffer) []Note {
	var notes []Note

	for true {
		var jumps uint16
		_ = binary.Read(buf, binary.LittleEndian, &jumps)

		if jumps == 0 {
			break
		}

		instrument, _ := buf.ReadByte()
		key, _ := buf.ReadByte()
		volume, _ := buf.ReadByte()
		pan, _ := buf.ReadByte()
		var pitch uint16
		_ = binary.Read(buf, binary.LittleEndian, &pitch)

		note := Note{
			instrument, key, volume, pan, pitch,
		}

		notes = append(notes, note)
	}

	return notes
}
