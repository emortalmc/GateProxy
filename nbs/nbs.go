package nbs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
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
	log.Println("ReadFile")
	file, err := os.ReadFile(path)
	if err != nil {
		return NBS{}, err
	}

	buf := bytes.NewBuffer(file)

	// First two bytes are always zero
	log.Println("First byte")
	byt, err := buf.ReadByte()
	if err != nil {
		return NBS{}, err
	}
	if byt != 0 {
		return NBS{}, fmt.Errorf("first byte is not zero, instead %d", byt)
	}
	log.Println("Second byte")
	byt, err = buf.ReadByte()
	if err != nil {
		return NBS{}, err
	}
	if byt != 0 {
		return NBS{}, fmt.Errorf("second byte is not zero, instead %d", byt)
	}

	nbs := NBS{}

	version, err := buf.ReadByte()
	log.Printf("Version: %d\n", version)
	if err != nil {
		return NBS{}, err
	}
	instruments, err := buf.ReadByte()
	if err != nil {
		return NBS{}, err
	}
	var length uint16
	err = binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return NBS{}, err
	}
	var layers uint16
	err = binary.Read(buf, binary.LittleEndian, &layers)
	if err != nil {
		return NBS{}, err
	}

	name, err := readString(buf)
	if err != nil {
		return NBS{}, err
	}
	author, err := readString(buf)
	if err != nil {
		return NBS{}, err
	}
	originalAuthor, err := readString(buf)
	if err != nil {
		return NBS{}, err
	}
	description, err := readString(buf)
	if err != nil {
		return NBS{}, err
	}

	var tps uint16
	err = binary.Read(buf, binary.LittleEndian, &tps)
	if err != nil {
		return NBS{}, err
	}

	autoSaving, err := buf.ReadByte()
	if err != nil {
		return NBS{}, err
	}
	autoSavingDuration, err := buf.ReadByte()
	if err != nil {
		return NBS{}, err
	}
	timeSignature, err := buf.ReadByte()
	if err != nil {
		return NBS{}, err
	}

	var minutesSpent int32
	err = binary.Read(buf, binary.LittleEndian, &minutesSpent)
	if err != nil {
		return NBS{}, err
	}
	var leftClicks int32
	err = binary.Read(buf, binary.LittleEndian, &leftClicks)
	if err != nil {
		return NBS{}, err
	}
	var rightClicks int32
	err = binary.Read(buf, binary.LittleEndian, &rightClicks)
	if err != nil {
		return NBS{}, err
	}
	var notesAdded int32
	err = binary.Read(buf, binary.LittleEndian, &notesAdded)
	if err != nil {
		return NBS{}, err
	}
	var notesRemoved int32
	err = binary.Read(buf, binary.LittleEndian, &notesRemoved)
	if err != nil {
		return NBS{}, err
	}

	midiName, err := readString(buf)
	if err != nil {
		return NBS{}, err
	}
	loop, err := buf.ReadByte()
	if err != nil {
		return NBS{}, err
	}
	maxLoops, err := buf.ReadByte()
	if err != nil {
		return NBS{}, err
	}

	var loopStart uint16
	err = binary.Read(buf, binary.LittleEndian, &loopStart)
	if err != nil {
		return NBS{}, err
	}

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

	nbs.Ticks, err = readNotes(buf, nbs)
	if err != nil {
		return NBS{}, err
	}

	fmt.Printf("\n\nRead %d ticks \n\n", len(nbs.Ticks))

	return nbs, nil
}

func readString(buf *bytes.Buffer) (string, error) {
	var lengthOfUtf int32
	err := binary.Read(buf, binary.LittleEndian, &lengthOfUtf)
	if err != nil {
		return "", err
	}
	utfBytes := make([]byte, lengthOfUtf)
	_, err = buf.Read(utfBytes)
	if err != nil {
		return "", err
	}

	return string(utfBytes), nil
}

func readNotes(buf *bytes.Buffer, nbs NBS) ([]Tick, error) {
	tick := -1
	ticks := make([]Tick, nbs.Length+1)

	for true {
		var jumps uint16
		err := binary.Read(buf, binary.LittleEndian, &jumps)
		if err != nil {
			return ticks, err
		}

		if jumps == 0 {
			break
		}

		tick += int(jumps)

		notes, err := readTickNoteLayers(buf)
		if err != nil {
			return ticks, err
		}

		ticks[tick] = Tick{
			Tick:  uint16(tick),
			Notes: notes,
		}
	}

	return ticks, nil
}

func readTickNoteLayers(buf *bytes.Buffer) ([]Note, error) {
	var notes []Note

	for true {
		var jumps uint16
		err := binary.Read(buf, binary.LittleEndian, &jumps)
		if err != nil {
			return notes, err
		}

		if jumps == 0 {
			break
		}

		instrument, err := buf.ReadByte()
		if err != nil {
			return notes, err
		}
		key, err := buf.ReadByte()
		if err != nil {
			return notes, err
		}
		volume, err := buf.ReadByte()
		if err != nil {
			return notes, err
		}
		pan, err := buf.ReadByte()
		if err != nil {
			return notes, err
		}
		var pitch uint16
		err = binary.Read(buf, binary.LittleEndian, &pitch)
		if err != nil {
			return notes, err
		}

		note := Note{
			instrument, key, volume, pan, pitch,
		}

		notes = append(notes, note)
	}

	return notes, nil
}
