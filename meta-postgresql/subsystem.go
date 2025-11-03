package metapostgresql

import (
	"fmt"
	"meta/engine"
	metaerror "meta/meta-error"
	metaflag "meta/meta-flag"
	"meta/meta-sql"
	"meta/subsystem"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Subsystem struct {
	subsystem.Subsystem
	GetConfig       func() map[string]*Config
	postgresqlDbs   map[string]*gorm.DB
	defaultSSLMode  string
	defaultTimeZone string
}

func GetSubsystem() *Subsystem {
	if thisSubsystem := engine.GetSubsystem[*Subsystem](); thisSubsystem != nil {
		return thisSubsystem.(*Subsystem)
	}
	return nil
}

func (s *Subsystem) GetName() string {
	return "PostgreSQL"
}

func (s *Subsystem) Start() error {
	config := s.GetConfig()
	if config == nil {
		return metaerror.New("postgres config is nil")
	}

	if s.defaultSSLMode == "" {
		s.defaultSSLMode = "disable"
	}
	if s.defaultTimeZone == "" {
		s.defaultTimeZone = "Asia/Shanghai"
	}

	if s.postgresqlDbs == nil {
		s.postgresqlDbs = make(map[string]*gorm.DB)
	}

	for key, cfg := range config {
		if cfg == nil {
			return metaerror.New("postgres config is nil")
		}
		if cfg.Host == "" {
			return metaerror.New("postgres host is empty")
		}
		if cfg.Port == 0 {
			cfg.Port = 5432
		}
		if cfg.Username == "" {
			return metaerror.New("postgres username is empty")
		}
		if cfg.Password == "" {
			return metaerror.New("postgres password is empty")
		}
		if cfg.Database == "" {
			return metaerror.New("postgres database is empty")
		}

		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
		)

		db, err := gorm.Open(
			postgres.Open(dsn), &gorm.Config{
				Logger: initGormLogger(),
			},
		)
		if err != nil {
			return metaerror.Wrap(err, fmt.Sprintf("failed to connect to PostgreSQL for key: %s", key))
		}

		sqlDB, err := db.DB()
		if err != nil {
			return metaerror.Wrap(err, "get underlying sql.DB failed")
		}
		// 设置连接池参数，可根据实际需求调整
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(time.Hour)

		s.postgresqlDbs[key] = db
	}

	return nil
}

func (s *Subsystem) Stop() error {
	var finalErr error
	for key, db := range s.postgresqlDbs {
		sqlDB, err := db.DB()
		if err != nil {
			finalErr = metaerror.Join(finalErr, metaerror.Wrap(err, fmt.Sprintf("[%s] get sql.DB failed", key)))
			continue
		}
		if err := sqlDB.Close(); err != nil {
			finalErr = metaerror.Join(finalErr, metaerror.Wrap(err, fmt.Sprintf("[%s] close DB failed", key)))
		}
	}
	return finalErr
}

func (s *Subsystem) GetClient(key string) *gorm.DB {
	if db, ok := s.postgresqlDbs[key]; ok {
		return db
	}
	return nil
}

func initGormLogger() gormlogger.Interface {
	loggerConfig := gormlogger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  gormlogger.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	}
	if metaflag.IsDebug() {
		loggerConfig.LogLevel = gormlogger.Info
	}
	return &metasql.Logger{
		Config: loggerConfig,
	}
}
