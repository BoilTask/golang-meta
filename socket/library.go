package socket

import (
	"fmt"
	googleProto "github.com/golang/protobuf/proto"
	"golang.org/x/exp/constraints"
	"log/slog"
	metaerror "meta/meta-error"
	"net"
)

func Connect(host string, port int32) (int32, error) {
	socketSubsystem := GetSubsystem()
	if socketSubsystem == nil {
		return -1, metaerror.New("socket subsystem not found")
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return -1, metaerror.Wrap(err, "error connecting to server")
	}

	socketIndex := socketSubsystem.indexGenerator.Next() // 获取下一个 ID
	socket := NewSocket(socketIndex, conn)               // 创建 Socket 实例
	socketSubsystem.socketsMutex.Lock()
	socketSubsystem.sockets[socketIndex] = socket // 存储 Socket 实例
	socketSubsystem.socketsMutex.Unlock()
	socket.Start(
		func() {
			socketSubsystem.socketsMutex.Lock()
			delete(socketSubsystem.sockets, socketIndex)
			socketSubsystem.socketsMutex.Unlock()
		},
	)

	slog.Info("Connected to server", "host", host, "port", port, "socketIndex", socketIndex, "Addr", conn.RemoteAddr())

	return socketIndex, nil
}

func SendMessage[T constraints.Integer](socketIndex int32, messageId T, proto googleProto.Message) (
	int32,
	error,
) {
	socketSubsystem := GetSubsystem()
	if socketSubsystem == nil {
		return -1, metaerror.New("socket subsystem not found")
	}
	socketSubsystem.socketsMutex.RLock()
	socket, exists := socketSubsystem.sockets[socketIndex]
	socketSubsystem.socketsMutex.RUnlock()

	if !exists {
		slog.Error("Socket not found", "socketIndex", socketIndex, "messageId", messageId)
		return -1, metaerror.New("socket not found: %d", socketIndex)
	}

	return socket.Send(int32(messageId), proto)
}

func GetMessageChannelBySocketIndex(index int32) string {
	return fmt.Sprintf("socket-%d", index)
}

func GetMessageChannelByMessageId[T constraints.Integer](id T) string {
	return fmt.Sprintf("id-%d", id)
}
