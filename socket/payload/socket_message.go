package payload

import (
	"github.com/golang/protobuf/proto"
	"meta/event"
	metaerror "meta/meta-error"
)

type SocketMessage struct {
	SocketIndex int32
	MessageId   int32
	ProtoByte   []byte
}

func ParseSocketMessage[T any, _ interface {
	*T
	proto.Message
}](payload []event.Payload) (*SocketMessage, *T, error) {
	if len(payload) != 1 {
		return nil, nil, metaerror.New("invalid payload length: %d", len(payload))
	}
	packet := payload[0].(*SocketMessage)
	var protoData T
	message, ok := interface{}(&protoData).(proto.Message)
	if !ok {
		return nil, nil, metaerror.New("type T %T does not implement proto.Message", protoData)
	}
	err := proto.Unmarshal(packet.ProtoByte, message)
	if err != nil {
		return nil, nil, err
	}
	return packet, &protoData, nil
}

type SocketMessageBuilder struct {
	socketMessage *SocketMessage
}

func NewSocketMessageBuilder() *SocketMessageBuilder {
	return &SocketMessageBuilder{socketMessage: &SocketMessage{}}
}

func (b *SocketMessageBuilder) SocketIndex(socketIndex int32) *SocketMessageBuilder {
	b.socketMessage.SocketIndex = socketIndex
	return b
}

func (b *SocketMessageBuilder) MessageId(messageId int32) *SocketMessageBuilder {
	b.socketMessage.MessageId = messageId
	return b
}

func (b *SocketMessageBuilder) ProtoByte(proto []byte) *SocketMessageBuilder {
	b.socketMessage.ProtoByte = proto
	return b
}

func (b *SocketMessageBuilder) Build() *SocketMessage {
	return b.socketMessage
}
