package metamysql

import (
	"fmt"
	"meta/engine"
	metaerror "meta/meta-error"
	metaflag "meta/meta-flag"
	"meta/meta-sql"
	"meta/subsystem"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Subsystem struct {
	subsystem.Subsystem
	GetConfig func() map[string]*Config
	mysqlDbs  map[string]*gorm.DB
}

func GetSubsystem() *Subsystem {
	if thisSubsystem := engine.GetSubsystem[*Subsystem](); thisSubsystem != nil {
		return thisSubsystem.(*Subsystem)
	}
	return nil
}

func (s *Subsystem) GetName() string {
	return "Mysql"
}

func (s *Subsystem) Start() error {
	config := s.GetConfig()
	if config == nil {
		return metaerror.New("mysql config is nil")
	}

	for key, config := range config {
		if config == nil {
			return metaerror.New("mysql config is nil")
		}
		if config.Host == "" {
			return metaerror.New("mysql uri is empty")
		}
		if config.Port == 0 {
			config.Port = 3306
		}
		if config.Username == "" {
			return metaerror.New("mysql username is empty")
		}
		if config.Password == "" {
			return metaerror.New("mysql password is empty")
		}
		if config.Database == "" {
			return metaerror.New("mysql database is empty")
		}

		mysqlDSN := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.Database,
		)

		loggerConfig := gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}
		if metaflag.IsDebug() {
			loggerConfig.LogLevel = gormlogger.Info
		}
		gormlogger.Default = &metasql.Logger{
			loggerConfig,
		}

		db, err := gorm.Open(mysql.Open(mysqlDSN))
		if err != nil {
			return metaerror.Wrap(err, "failed to connect to MySQL")
		}
		sqlDB, err := db.DB()
		if err != nil {
			return metaerror.Wrap(err, "failed to get sql.DB")
		}
		// 设置最大连接数
		sqlDB.SetMaxOpenConns(50)
		// 设置最大空闲连接数
		sqlDB.SetMaxIdleConns(10)
		// 设置空闲连接的最大生命周期（比如 1 分钟）
		sqlDB.SetConnMaxIdleTime(1 * time.Minute)
		// 设置连接的最大生命周期（比如 5 分钟，防止被服务端 kill 掉）
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
		if s.mysqlDbs == nil {
			s.mysqlDbs = make(map[string]*gorm.DB)
		}
		s.mysqlDbs[key] = db
	}
	return nil
}

func (s *Subsystem) Stop() error {
	var finalErr error
	if s.mysqlDbs != nil {
		for _, db := range s.mysqlDbs {
			sqlDB, err := db.DB()
			if err != nil {
				finalErr = metaerror.Join(finalErr, metaerror.Wrap(err, "get underlying sql.DB failed"))
				continue
			}
			if err := sqlDB.Close(); err != nil {
				finalErr = metaerror.Join(finalErr, metaerror.Wrap(err, "gorm close failed"))
			}
		}
	}
	return finalErr
}

func (s *Subsystem) GetClient(key string) *gorm.DB {
	if s.mysqlDbs == nil {
		return nil
	}
	if db, ok := s.mysqlDbs[key]; ok {
		return db
	}
	return nil
}
