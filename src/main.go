package main

import (
	"runtime"
	"sync"

	"app/api"
	"app/bellhop"
	"app/db"
	"app/metrics"
	"app/models"
	"app/redis"
	"app/service"

	"github.com/joho/godotenv"
	"gitlab.online-fx.com/go-packages/apiresponse"
	"gitlab.online-fx.com/go-packages/envchecker"
	"gitlab.online-fx.com/go-packages/logger"
	"gitlab.online-fx.com/go-packages/prometheus"
)

// init is invoked before main().
func init() {
	requiredVariables := models.GetRequiredVariables()

	// Проверяем что уже установлены все необходимые для работы сервиса переменные окружения
	if errKubernetes := envchecker.CheckEnvironments(requiredVariables); errKubernetes != nil {
		logger.Errorf(envchecker.ErrorKubernetesMessage+" Error: %v", errKubernetes)

		// loads values from config.conf into the system
		if err := godotenv.Load("config/config.conf"); err != nil {
			logger.Fatal("file config.conf not found")
		}

		// loads values from env.conf into the system
		if err := godotenv.Load("config/env.conf"); err != nil {
			logger.Error("file env.conf not found")
		}

		// Проверяем что в конфигурационных файлах установлены все необходимые для работы сервиса переменные окружения
		if errConfigs := envchecker.CheckEnvironments(requiredVariables); errConfigs != nil {
			logger.Fatalf(envchecker.ErrorConfigFilesMessage+" Error: %v", errConfigs)
		}
	}
}

func main() {
	// Загружаем конфигурацию сервиса
	config, err := models.LoadConfig()
	if err != nil {
		logger.Fatalf("Config error: %v", err)
	}

	// Подключение к БД
	_, err = db.Init(config)
	if err != nil {
		logger.Fatalf("Database connection error: %v", err)
	}

	// Запуск миграций
	db.RunMigrations()
	logger.Info("DB migrated")

	// Инициализация клиента Redis
	redisClient := redis.New(config)
	logger.Info("Redis initialized")

	// Инициализация пушера в Bellhop
	bellhop.InitConfig(config)

	// Инициализация метрик
	metrics.InitMetrics(config.MetricSettings)

	// Инициализация сервиса и памяти
	service.Load(config, redisClient)
	service.Restore()

	// Инициализация API route
	apiresponse.InitRoutes(api.Routes)
	logger.Info("Api routes initialized")

	// Инициализация метрик
	prometheus.InitMetricsRoute(config.Prometheus)
	go prometheus.InitServer(config.Prometheus)

	var wg sync.WaitGroup
	wg.Add(1)

	logger.Info("Service started")
	logger.Infof("Go version: %s", runtime.Version())
	wg.Wait()
}
