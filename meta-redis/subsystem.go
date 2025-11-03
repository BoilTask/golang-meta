package metaredis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"meta/engine"
	metaerror "meta/meta-error"
	"meta/subsystem"
)

type Subsystem struct {
	subsystem.Subsystem
	GetConfig func() *Config
	client    *redis.Client
}

func GetSubsystem() *Subsystem {
	if thisSubsystem := engine.GetSubsystem[*Subsystem](); thisSubsystem != nil {
		return thisSubsystem.(*Subsystem)
	}
	return nil
}

func (redisSubsystem *Subsystem) GetName() string {
	return "Redis"
}

func (redisSubsystem *Subsystem) Start() error {
	config := redisSubsystem.GetConfig()
	if config == nil {
		return metaerror.New("redis config is nil")
	}
	redisSubsystem.client = redis.NewClient(
		&redis.Options{
			Addr:     config.Addr,     // Redis 服务器地址
			Password: config.Password, // 如果没有密码则留空
		},
	)
	ctx := context.Background()
	if _, err := redisSubsystem.client.Ping(ctx).Result(); err != nil {
		return err
	}
	return nil
}

func (redisSubsystem *Subsystem) Stop() error {
	if redisSubsystem.client != nil {
		if err := redisSubsystem.client.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (redisSubsystem *Subsystem) GetClient() *redis.Client {
	return redisSubsystem.client
}
