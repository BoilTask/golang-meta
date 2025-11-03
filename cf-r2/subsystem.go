package cfr2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"meta/engine"
	metaerror "meta/meta-error"
	"meta/subsystem"
	"net/http"
)

type Subsystem struct {
	subsystem.Subsystem
	GetConfig func() map[string]*Config
	Clients   map[string]*s3.S3
}

func GetSubsystem() *Subsystem {
	if thisSubsystem := engine.GetSubsystem[*Subsystem](); thisSubsystem != nil {
		return thisSubsystem.(*Subsystem)
	}
	return nil
}

func (s *Subsystem) GetName() string {
	return "CfR2"
}

func (s *Subsystem) Start() error {
	config := s.GetConfig()
	if config == nil {
		return metaerror.New("cf-r2 config is nil")
	}

	s.Clients = make(map[string]*s3.S3)

	for key, config := range config {
		if config == nil {
			return metaerror.New("cf-r2 config is nil")
		}
		r2Session, err := session.NewSession(&aws.Config{
			Region:           aws.String("auto"),     // R2一般写 auto
			Endpoint:         aws.String(config.Url), // 替换成你的 R2 Endpoint
			S3ForcePathStyle: aws.Bool(true),         // R2要求这个必须 true
			Credentials: credentials.NewStaticCredentials(config.Key,
				config.Secret,
				config.Token),
		})
		if err != nil {
			return metaerror.Wrap(err, "create session failed, key:%s", key)
		}
		s3Client := s3.New(r2Session)
		if s3Client == nil {
			return metaerror.New("create s3 client failed, key:%s", key)
		}
		s.Clients[key] = s3Client
	}
	return nil
}

func (s *Subsystem) Stop() error {
	var finalErr error
	for key, client := range s.Clients {
		if tr, ok := client.Config.HTTPClient.Transport.(*http.Transport); ok {
			tr.CloseIdleConnections()
		}
		delete(s.Clients, key)
	}
	return finalErr
}

func (s *Subsystem) GetClient(key string) *s3.S3 {
	client, ok := s.Clients[key]
	if ok {
		return client
	}
	return nil
}
