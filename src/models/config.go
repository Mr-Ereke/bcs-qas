package models

import (
	"os"
	"strconv"

	"gitlab.online-fx.com/go-packages/gormdb"
	"gitlab.online-fx.com/go-packages/prometheus"
	"gorm.io/gorm/schema"
)

const ServiceDB = "ServiceDB"

type Config struct {
	Database       gormdb.Database
	Redis          RedisConfig
	Bellhop        Bellhop
	Settings       Settings
	MetricSettings MetricSettings
	Prometheus     prometheus.Prometheus
	LogLevel       logLevelType
	Timezone       string
}

type Database struct {
	Address        string
	Port           int
	User           string
	Password       string
	DB             string
	MaxConnections int
}

type RedisConfig struct {
	Address     string
	Port        int
	Password    string
	DB          int
	PingTimeout int
}

type Bellhop struct {
	Host     string
	Timeout  int
	PushName string
}

type Settings struct {
	MaxAlertsCustomerInstrument int
	MaxAlertsCustomer           int
}

type MetricSettings struct {
	LongPushTime int64
}

type logLevelType int

const (
	LevelPanic logLevelType = iota
	LevelFatal
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
)

func LoadConfig() (*Config, error) {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	if err != nil {
		redisDB = 0
	}

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return nil, err
	}

	redisPingTimeout, err := strconv.Atoi(os.Getenv("REDIS_PING_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	bellhopTimeout, err := strconv.Atoi(os.Getenv("BELLHOP_TIMEOUT"))
	if err != nil {
		return nil, err
	}

	customerInstrumentAlerts, err := strconv.Atoi(os.Getenv("MAX_COUNT_ALERTS_CUSTOMER_INSTRUMENT"))
	if err != nil {
		return nil, err
	}

	customerAlerts, err := strconv.Atoi(os.Getenv("MAX_COUNT_ALERTS_CUSTOMER"))
	if err != nil {
		return nil, err
	}

	metricsLongPushTime, err := strconv.Atoi(os.Getenv("METRICS_LONG_PUSH_TIME"))
	if err != nil {
		return nil, err
	}

	logLevel, _ := strconv.Atoi(os.Getenv("LOG_LEVEL"))

	return &Config{
		Prometheus: prometheus.Prometheus{
			Port: os.Getenv("SERVICE_PORT"),
		},
		Database: gormdb.Database{
			Address:  os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASS"),
			DB:       os.Getenv("DATABASE_NAME"),
		},
		Redis: RedisConfig{
			Address:     os.Getenv("REDIS_ADDRESS"),
			Port:        redisPort,
			Password:    os.Getenv("REDIS_PASSWORD"),
			DB:          redisDB,
			PingTimeout: redisPingTimeout,
		},
		Bellhop: Bellhop{
			Host:     os.Getenv("BELLHOP_HOST"),
			Timeout:  bellhopTimeout,
			PushName: os.Getenv("BELLHOP_PUSH_NAME"),
		},
		Settings: Settings{
			MaxAlertsCustomerInstrument: customerInstrumentAlerts,
			MaxAlertsCustomer:           customerAlerts,
		},
		MetricSettings: MetricSettings{
			LongPushTime: int64(metricsLongPushTime),
		},
		LogLevel: logLevelType(logLevel),
		Timezone: os.Getenv("TIMEZONE"),
	}, nil
}

func GetRequiredVariables() []string {
	return []string{
		// Обязательные переменные окружения для сервиса
		"SERVICE_PORT",

		// Подключение к БД
		"DATABASE_HOST", "DATABASE_PORT", "DATABASE_USER", "DATABASE_PASS", "DATABASE_NAME",

		// Обязательные переменные окружения для подключения к Redis
		"REDIS_ADDRESS", "REDIS_PORT", "REDIS_DB_NUMBER", "REDIS_PING_TIMEOUT",

		// Обязательные переменные окружения для отправки пушей в Bellhop
		"BELLHOP_HOST", "BELLHOP_TIMEOUT", "BELLHOP_PUSH_NAME",

		// Обязательные переменные окружения валидации количества алертов
		"MAX_COUNT_ALERTS_CUSTOMER_INSTRUMENT", "MAX_COUNT_ALERTS_CUSTOMER",

		// Обязательные переменные окружения для метрик
		"METRICS_LONG_PUSH_TIME",

		// Временная зона
		"TIMEZONE",

		// Уровень логирования
		"LOG_LEVEL",
	}
}

func GetModels() []schema.Tabler {
	return []schema.Tabler{
		Subscription{},
	}
}
