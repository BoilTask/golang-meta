package socket

import (
	"context"
	"log/slog"
	"meta/event"
	metaerror "meta/meta-error"
	"meta/metaroutine"
	"meta/network"
	socketEvent "meta/socket/event"
	socketPayload "meta/socket/payload"
	"net"
	"sync"

	googleProto "github.com/golang/protobuf/proto"
)

// Socket 封装 net.Conn 和消息通道
type Socket struct {
	socketIndex    int32
	conn           net.Conn
	dataChan       chan []byte
	sendBuffers    net.Buffers
	receiveBuffers net.Buffers
}

// NewSocket 创建新的 Socket 实例
func NewSocket(socketIndex int32, conn net.Conn) *Socket {
	return &Socket{
		socketIndex:    socketIndex,
		conn:           conn,
		dataChan:       make(chan []byte, 100),
		sendBuffers:    make(net.Buffers, 10),
		receiveBuffers: make(net.Buffers, 10),
	}
}

// Start 启动消息处理
func (s *Socket) Start(callback func()) {
	metaroutine.SafeGoWithRestart(
		"Socket start",
		func() error {
			return s.processSocket(callback)
		},
	)
}

// Send 发送消息到连接
func (s *Socket) Send(messageId int32, proto googleProto.Message) (int32, error) {
	metaroutine.SafeGo(
		"Socket Send",
		func() error {
			var protoBytes []byte
			var err error
			if proto != nil {
				protoBytes, err = googleProto.Marshal(proto)
				if err != nil {
					return metaerror.Wrap(err, "error marshaling message")
				}
			}
			packet := network.ConvertPacket(nil, -1, -1, messageId, int32(len(protoBytes)), protoBytes)
			networkBytes, packageErr := network.MakeBytes(packet)
			if packageErr != nil {
				return metaerror.Wrap(packageErr, "error making package")
			}
			slog.Info("Message sent", "messageId", messageId, "dataSize", len(networkBytes))
			s.dataChan <- networkBytes
			return nil
		},
	)
	return -1, nil
}

func (s *Socket) processSocket(callback func()) error {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)

	metaroutine.SafeGoWithRestart(
		"SocketReceive",
		func() error {
			s.handleReceiveMessages(ctx, &wg)
			cancel()
			return nil
		},
	)

	metaroutine.SafeGoWithRestart(
		"SocketSend",
		func() error {
			s.handleSendMessages(ctx, &wg)
			cancel()
			return nil
		},
	)

	channels := []string{
		GetMessageChannelBySocketIndex(s.socketIndex),
	}
	payload := socketPayload.NewSocketBuilder().
		SocketIndex(s.socketIndex).
		Build()
	event.InvokeChannel[socketEvent.SocketConnected](&channels, payload)

	wg.Wait()

	slog.Info("socket End", "index", s.socketIndex)

	callback()

	return nil
}

func (s *Socket) handleReceiveMessages(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		err := s.Close()
		if err != nil {
			return
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			loadPackage, err := network.LoadPackage(s.conn)
			if err != nil {
				return
			}
			packets, err := network.ParsePacket(loadPackage)
			if err != nil {
				continue
			}
			for _, packet := range packets {
				channels := []string{
					GetMessageChannelBySocketIndex(s.socketIndex),
					GetMessageChannelByMessageId(packet.MessageId),
				}
				payload := socketPayload.NewSocketMessageBuilder().
					SocketIndex(s.socketIndex).
					MessageId(packet.MessageId).
					ProtoByte(packet.ProtoData).
					Build()
				event.InvokeChannel[socketEvent.SocketMessage](&channels, payload)
			}
		}
	}
}

func (s *Socket) handleSendMessages(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range s.dataChan {
		select {
		case <-ctx.Done():
			return
		default:
			s.sendBuffers = append(s.sendBuffers, data)
			if _, err := s.sendBuffers.WriteTo(s.conn); err != nil {
				return
			}
		}
	}
}

// Close 关闭连接和消息通道
func (s *Socket) Close() error {
	defer func() {
		channels := []string{
			GetMessageChannelBySocketIndex(s.socketIndex),
		}
		payload := socketPayload.NewSocketBuilder().
			SocketIndex(s.socketIndex).
			Build()
		event.InvokeChannel[socketEvent.SocketDisconnected](&channels, payload)
	}()
	if s.conn != nil {
		close(s.dataChan)
	}
	err := s.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// IsConnected 检查连接状态
func (s *Socket) IsConnected() bool {
	return s.conn != nil
}
