package internal

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/phantomnat/wtf-dial/internal/wtf"
)

//go:generate protoc --gogo_out=. internal.proto

// MarshalDial encodes a dial to binary format.
func MarshalDial(d *wtf.Dial) ([]byte, error) {
	return proto.Marshal(&Dial{
		ID:      int64(d.ID),
		UserID:  int64(d.UserID),
		Name:    d.Name,
		Level:   d.Level,
		ModTime: d.ModTime.UnixNano(),
	})
}

// UnmarshalDial decodes a dial from a binary data.
func UnmarshalDial(data []byte, d *wtf.Dial) error {
	var pb Dial
	if err := proto.Unmarshal(data, &pb); err != nil {
		return err
	}

	d.ID = wtf.DialID(pb.GetID())
	d.UserID = wtf.UserID(pb.GetUserID())
	d.Name = pb.GetName()
	d.Level = pb.GetLevel()
	d.ModTime = time.Unix(0, pb.GetModTime()).UTC()

	return nil
}
