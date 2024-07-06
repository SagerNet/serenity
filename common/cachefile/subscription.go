package cachefile

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json"
	"github.com/sagernet/sing/common/varbin"
)

type Subscription struct {
	Content     []option.Outbound
	LastUpdated time.Time
	LastEtag    string
}

func (c *Subscription) MarshalBinary() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteByte(1)
	content, err := json.Marshal(c.Content)
	if err != nil {
		return nil, err
	}
	_, err = varbin.WriteUvarint(&buffer, uint64(len(content)))
	if err != nil {
		return nil, err
	}
	_, err = buffer.Write(content)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buffer, binary.BigEndian, c.LastUpdated.Unix())
	if err != nil {
		return nil, err
	}
	err = varbin.Write(&buffer, binary.BigEndian, c.LastEtag)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (c *Subscription) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)
	version, err := reader.ReadByte()
	if err != nil {
		return err
	}
	_ = version
	contentLength, err := binary.ReadUvarint(reader)
	if err != nil {
		return err
	}
	content := make([]byte, contentLength)
	_, err = reader.Read(content)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &c.Content)
	if err != nil {
		return err
	}
	var lastUpdatedUnix int64
	err = binary.Read(reader, binary.BigEndian, &lastUpdatedUnix)
	if err != nil {
		return err
	}
	c.LastUpdated = time.Unix(lastUpdatedUnix, 0)
	err = varbin.Read(reader, binary.BigEndian, &c.LastEtag)
	if err != nil {
		return err
	}
	return nil
}
