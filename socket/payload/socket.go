package payload

type Socket struct {
	SocketIndex int32
}

type SocketBuilder struct {
	socket *Socket
}

func NewSocketBuilder() *SocketBuilder {
	return &SocketBuilder{socket: &Socket{}}
}

func (b *SocketBuilder) SocketIndex(socketIndex int32) *SocketBuilder {
	b.socket.SocketIndex = socketIndex
	return b
}

func (b *SocketBuilder) Build() *Socket {
	return b.socket
}
