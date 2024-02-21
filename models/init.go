package models

import (
	"context"
	"time"

	"github.com/lazyliqiquan/help_rookie/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	logger *zap.Logger
	DB     *gorm.DB
	RDB    *redis.Client
)

func Init(loggerInstance *zap.Logger, config *config.WebConfig) {
	logger = loggerInstance
	var err error
	mysqlDsn := "root:" + config.MysqlPassword + "@tcp(" +
		config.MysqlPath + ":" + config.MysqlPort
	if config.Debug {
		mysqlDsn = "root:" + config.MysqlPassword + "@tcp(" +
			config.DebugMysqlPath + ":" + config.DebugMysqlPort
	}
	mysqlDsn += ")/?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Connect mysql fail : ", zap.Error(err))
	}
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Fatal("Give sql.DB fail : ", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	err = DB.Exec("CREATE DATABASE IF NOT EXISTS help_cookie").Error
	if err != nil {
		logger.Fatal("Create database help_cookie fail : ", zap.Error(err))
	}
	err = DB.Exec("USE help_cookie").Error
	if err != nil {
		logger.Fatal("Unable to use the database help_cookie : ", zap.Error(err))
	}
	err = DB.AutoMigrate(&Users{}, &SeekHelps{}, &LendHands{}, &Comments{})
	if err != nil {
		logger.Fatal("Create tables fail : ", zap.Error(err))
	}
	logger.Sugar().Infoln("Help-cookie mysql init succeed !")

	redisDsn := config.RedisPath + ":" + config.RedisPort
	if config.Debug {
		redisDsn = config.DebugRedisPath + ":" + config.DebugRedisPort
	}
	RDB = redis.NewClient(&redis.Options{
		Addr:     redisDsn,
		Password: "",
	})
	WebGlobalParams := config.RedisInit()
	for k, v := range WebGlobalParams {
		if err := RDB.Set(context.Background(), k, v, time.Duration(0)).Err(); err != nil {
			logger.Fatal("Set web global params fail : ", zap.Error(err))
		}
	}
	logger.Sugar().Infoln("Help-cookie redis init succeed !")
}
