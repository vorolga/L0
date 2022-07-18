package main

import (
	"L0/connecting"
	"L0/internal/handlers"
	"L0/internal/repository"
	"L0/internal/usecase"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	echoServer := echo.New()

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, err := config.Build()
	if err != nil {
		log.Fatal("zap logger build error")
	}
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		err = prLogger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(prLogger)

	if err = godotenv.Load(".env"); err != nil {
		logger.Fatal("error while loading environment ", err)
	}

	postgres, err := connecting.NewPostgresDBComposite()
	if err != nil {
		logger.Fatal("postgres composite failed ", err)
	}

	redis, err := connecting.NewRedisComposite()
	if err != nil {
		logger.Fatal("redis composite failed", err)
	}
	rep := repository.NewRepository(postgres.DB, redis.Redis)

	useCase, err := usecase.NewUsecase(rep)
	if err != nil {
		logger.Error("usecase failed ", err)
	}

	handler := handlers.NewHandler(useCase, logger)
	handler.Register(echoServer)
	if err != nil {
		logger.Error("handler failed ", err)
	}

	stanComposite, err := connecting.NewStanComposite()
	if err != nil {
		logger.Fatal("stan composite failed ", err)
	}

	_, err = stanComposite.SC.Subscribe(os.Getenv("NATS_CHANNEL"), handler.ReceiveOrder, stan.DurableName(os.Getenv("DURABLE")),
		stan.MaxInflight(25),
		stan.SetManualAckMode(),
	)
	if err != nil {
		logger.Fatal("subscribe failed ", err)
	}

	echoServer.Logger.Fatal(echoServer.Start(":8080"))
}
