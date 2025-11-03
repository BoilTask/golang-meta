package network

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Packet struct {
	UniqueIdSize int32
	UniqueId     []byte
	RequestId    int16
	ResponseId   int16
	MessageId    int32
	ProtoSize    int32
	ProtoData    []byte
}

func (p *Packet) MakeBytes(buffer *bytes.Buffer) error {
	err := binary.Write(buffer, binary.BigEndian, p.UniqueIdSize)
	if err != nil {
		return err
	}
	buffer.Write(p.UniqueId)
	err = binary.Write(buffer, binary.BigEndian, p.RequestId)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, p.ResponseId)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, p.MessageId)
	if err != nil {
		return err
	}
	err = binary.Write(buffer, binary.BigEndian, p.ProtoSize)
	if err != nil {
		return err
	}
	buffer.Write(p.ProtoData)
	return nil
}

func ParsePacket(p *Package) ([]*Packet, error) {

	var packets []*Packet
	for {

		buf := bytes.NewReader(p.Data)
		packet := &Packet{}

		if err := binary.Read(buf, binary.BigEndian, &packet.UniqueIdSize); err != nil {
			return nil, err
		}

		if packet.UniqueIdSize > 0 {
			packet.UniqueId = make([]byte, packet.UniqueIdSize)
			if _, err := io.ReadFull(buf, packet.UniqueId); err != nil {
				return nil, err
			}
		}

		if err := binary.Read(buf, binary.BigEndian, &packet.RequestId); err != nil {
			return nil, err
		}
		if err := binary.Read(buf, binary.BigEndian, &packet.ResponseId); err != nil {
			return nil, err
		}
		if err := binary.Read(buf, binary.BigEndian, &packet.MessageId); err != nil {
			return nil, err
		}

		if err := binary.Read(buf, binary.BigEndian, &packet.ProtoSize); err != nil {
			return nil, err
		}

		if packet.ProtoSize > 0 {
			packet.ProtoData = make([]byte, packet.ProtoSize)
			if _, err := io.ReadFull(buf, packet.ProtoData); err != nil {
				return nil, err
			}
		}

		packets = append(packets, packet)

		if buf.Len() == 0 {
			break
		}
	}

	return packets, nil
}

func ConvertPacket(
	uniqueId *string,
	requestId int16,
	responseId int16,
	messageId int32,
	protoSize int32,
	protoData []byte,
) *Packet {
	packet := &Packet{
		RequestId:  requestId,
		ResponseId: responseId,
		MessageId:  messageId,
		ProtoSize:  protoSize,
		ProtoData:  protoData,
	}
	if uniqueId != nil {
		packet.UniqueIdSize = int32(len(*uniqueId))
		packet.UniqueId = []byte(*uniqueId)
	} else {
		packet.UniqueIdSize = 0
		packet.UniqueId = []byte{}
	}
	return packet
}
