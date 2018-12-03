package proto

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/mafanr/g"
	"go.uber.org/zap"
)

// ApplicationStreamCreate ...
type ApplicationStreamCreate struct {
	Type      int16
	ChannelID int
	Length    int
	Payload   []byte
}

// NewApplicationStreamCreate ...
func NewApplicationStreamCreate() *ApplicationStreamCreate {
	return &ApplicationStreamCreate{
		Type: APPLICATION_STREAM_CREATE,
	}
}

// Decode ...
func (a *ApplicationStreamCreate) Decode(conn net.Conn, reader io.Reader) error {
	buf := make([]byte, 8)
	if _, err := io.ReadFull(reader, buf); err != nil {
		g.L.Warn("ApplicationRequest Decode", zap.String("error", err.Error()))
		return err
	}
	a.ChannelID = int(binary.BigEndian.Uint32(buf[:4]))
	a.Length = int(binary.BigEndian.Uint32(buf[4:8]))
	a.Payload = make([]byte, a.Length)

	if _, err := io.ReadFull(reader, a.Payload); err != nil {
		return err
	}
	//log.Println(string(a.Payload))
	return nil
}

// Encode ...
func (a *ApplicationStreamCreate) Encode() ([]byte, error) {
	return nil, nil
}

// GetPacketType ...
func (a *ApplicationStreamCreate) GetPacketType() int16 {
	return a.Type
}

// GetPayload ...
func (a *ApplicationStreamCreate) GetPayload() []byte {
	return nil
}

// GetRequestID ...
func (a *ApplicationStreamCreate) GetRequestID() int {
	return 0
}
