package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"log"
	math_rand "math/rand"

	"github.com/lazyliqiquan/help_rookie/config"
	"github.com/lazyliqiquan/help_rookie/middlewares"
	"github.com/lazyliqiquan/help_rookie/models"
	"github.com/lazyliqiquan/help_rookie/router"
	"github.com/lazyliqiquan/help_rookie/service"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	var err error
	logger, err = zap.NewProduction()
	defer logger.Sync()
	if err != nil {
		log.Fatalln("Init logger fail :", err)
	}
	models.Init(logger, config.Config)
	middlewares.Init(logger.Sugar())
	service.Init(logger.Sugar())
	initRand()
	r := router.Router(config.Config)
	if config.Config.Debug {
		r.Run(config.Config.DebugWebPath)
	} else {
		r.Run(config.Config.WebPath)
	}
}

func initRand() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		logger.Fatal("Random generator init failed : ", zap.Error(err))
	}
	sd := int64(binary.LittleEndian.Uint64(b[:]))
	logger.Sugar().Infof("random seed : %d ", sd)
	math_rand.Seed(sd)
}
