package socket

import (
	"fmt"
	"log/slog"
	"meta/engine"
	"meta/generator"
	"meta/subsystem"
	"net"
	"sync"
)

type Subsystem struct {
	subsystem.Subsystem
	GetPort        func() int32
	indexGenerator *generator.IncreaseGenerator[int32]
	sockets        map[int32]*Socket
	socketsMutex   sync.RWMutex
}

func GetSubsystem() *Subsystem {
	if thisSubsystem := engine.GetSubsystem[*Subsystem](); thisSubsystem != nil {
		return thisSubsystem.(*Subsystem)
	}
	return nil
}

func (socketSubsystem *Subsystem) GetName() string {
	return "Socket"
}

func (socketSubsystem *Subsystem) Init() error {
	socketSubsystem.sockets = map[int32]*Socket{}
	socketSubsystem.indexGenerator = generator.NewIncreaseGenerator[int32](0, 1)
	return nil
}

func (socketSubsystem *Subsystem) Start() error {
	if socketSubsystem.GetPort != nil {
		go socketSubsystem.startSubsystem()
	}
	return nil
}

func (socketSubsystem *Subsystem) startSubsystem() {

	port := socketSubsystem.GetPort()

	socketListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		slog.Error("starting server error", "err", err)
		return
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			slog.Error("Error closing listener", "err", err)
		}
	}(socketListener)

	slog.Info("Socket server is listening", "port", port)

	for {
		// 接受新的连接
		conn, err := socketListener.Accept()
		if err != nil {
			slog.Info("Error accepting connection:", err)
			continue
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

		slog.Info("New connection", "socketIndex", socketIndex, "Addr", conn.RemoteAddr())
	}
}
