package rpc

import (
	"fmt"
	"log/slog"
	"meta/engine"
	metaerror "meta/meta-error"
	"meta/metaroutine"
	"meta/subsystem"
	"net"

	"google.golang.org/grpc"
)

type Subsystem struct {
	subsystem.Subsystem
	GetPort func() int32
}

func GetSubsystem() *Subsystem {
	if thisSubsystem := engine.GetSubsystem[*Subsystem](); thisSubsystem != nil {
		return thisSubsystem.(*Subsystem)
	}
	return nil
}

func (s *Subsystem) GetName() string {
	return "Rpc"
}

func (s *Subsystem) Start() error {
	if s.GetPort != nil {
		metaroutine.SafeGoWithRestart(
			"Rpc start",
			s.startSubsystem,
		)
	}
	return nil
}

func (s *Subsystem) startSubsystem() error {
	port := s.GetPort()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return metaerror.Wrap(err, "failed to listen, port:%d", port)
	}
	defer func(lis net.Listener) {
		err := lis.Close()
		if err != nil {
			slog.Error("Error closing listener", "err", err)
		}
	}(lis)

	slog.Info("Rpc server is listening", "port", port)

	grpcServer := grpc.NewServer()
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
