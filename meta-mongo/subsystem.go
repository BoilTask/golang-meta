package metamongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"meta/engine"
	metaerror "meta/meta-error"
	metapanic "meta/meta-panic"
	"meta/subsystem"
	"sync"
	"time"
)

type Subsystem struct {
	subsystem.Subsystem
	GetConfig func() *Config
	client    *mongo.Client

	cmdMap sync.Map // 存储 RequestID -> command 文本
}

func GetSubsystem() *Subsystem {
	if thisSubsystem := engine.GetSubsystem[*Subsystem](); thisSubsystem != nil {
		return thisSubsystem.(*Subsystem)
	}
	return nil
}

func (s *Subsystem) GetName() string {
	return "Mongo"
}

func (s *Subsystem) Start() error {
	config := s.GetConfig()
	if config == nil {
		return metaerror.New("mongo config is nil")
	}
	uri := config.Uri
	var err error

	// 构造 CommandMonitor
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			s.cmdMap.Store(evt.RequestID, fmt.Sprint(evt.Command))
		},
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			// 从 map 里取出语句
			val, ok := s.cmdMap.Load(evt.RequestID)
			if ok {
				//slog.Info(
				//	"[Mongo] Succeeded",
				//	"command_name",
				//	evt.CommandName,
				//	"duration",
				//	fmt.Sprintf("%v", evt.Duration),
				//	"query",
				//	val,
				//)
				if evt.Duration > time.Millisecond*200 {
					slog.Warn(
						"Mongo command too slow",
						"command_name", evt.CommandName,
						"duration", fmt.Sprintf("%v", evt.Duration),
						"query", val,
					)
				}
				s.cmdMap.Delete(evt.RequestID) // 释放内存
			} else {
				metapanic.ProcessError(
					metaerror.New(
						"Mongo command succeeded but not found in map, command_name:%s, duration:%s",
						evt.CommandName,
						fmt.Sprintf("%v", evt.Duration),
					),
				)
			}
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			slog.Warn(
				"Mongo command failed",
				"command_name", evt.CommandName,
				"duration", fmt.Sprintf("%v", evt.Duration),
				"error", evt.Failure,
			)
		},
	}

	mongoOptions := options.Client().
		ApplyURI(uri).
		SetMonitor(monitor)

	mongoOptions.Auth = &options.Credential{
		Username: config.Username,
		Password: config.Password,
	}

	s.client, err = mongo.Connect(context.TODO(), mongoOptions)
	if err != nil {
		return err
	}
	return nil
}

func (s *Subsystem) Stop() error {
	if s.client != nil {
		if err := s.client.Disconnect(context.TODO()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Subsystem) GetClient() *mongo.Client {
	return s.client
}
