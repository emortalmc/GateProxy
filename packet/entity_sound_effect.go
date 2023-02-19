package packet

import (
	"go.minekube.com/gate/pkg/edition/java/proto/util"
	"go.minekube.com/gate/pkg/gate/proto"
	"io"
)

type EntitySoundEffect struct {
	SoundID       int
	SoundCategory int
	EntityID      int
	Volume        float32
	Pitch         float32
	Seed          int64
}

func (k *EntitySoundEffect) Encode(c *proto.PacketContext, wr io.Writer) error {
	err := util.WriteVarInt(wr, k.SoundID+1)
	if err != nil {
		return err
	}
	err = util.WriteVarInt(wr, k.SoundCategory)
	if err != nil {
		return err
	}
	err = util.WriteVarInt(wr, k.EntityID)
	if err != nil {
		return err
	}
	err = util.WriteFloat32(wr, k.Volume)
	if err != nil {
		return err
	}
	err = util.WriteFloat32(wr, k.Pitch)
	if err != nil {
		return err
	}
	err = util.WriteInt64(wr, k.Seed)
	if err != nil {
		return err
	}

	return nil
}

func (k *EntitySoundEffect) Decode(c *proto.PacketContext, rd io.Reader) (err error) {
	k.SoundID, _ = util.ReadVarInt(rd)
	k.SoundCategory, _ = util.ReadVarInt(rd)
	k.EntityID, _ = util.ReadVarInt(rd)
	k.Volume, _ = util.ReadFloat32(rd)
	k.Pitch, _ = util.ReadFloat32(rd)
	k.Seed, _ = util.ReadInt64(rd)

	return
}

var _ proto.Packet = (*EntitySoundEffect)(nil)
