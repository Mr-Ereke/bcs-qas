package metrics

import (
	"time"

	"app/models"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	pushTimeLatency     = "quotes_alerting_service_push_time_latency"
	pushTimeLatencyHelp = "Count of alerts exceeded sending time"
	lastPushTime        = "quotes_alerting_service_last_push_time"
	lastPushTimeHelp    = "Last time in milliseconds from alert before send push to Bellhop"
	bellhopLatency      = "quotes_alerting_service_bellhop_latency"
	bellhopLatencyHelp  = "Latency duration from send request to Bellhop before get response"
)

var (
	settings             models.MetricSettings
	longPushAlertCount   prometheus.Gauge
	lastPushTimeDuration prometheus.Gauge
	bellhopDuration      prometheus.Gauge
)

func InitMetrics(config models.MetricSettings) {
	settings = config
	longPushAlertCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: pushTimeLatency,
		Help: pushTimeLatencyHelp,
	})
	lastPushTimeDuration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: lastPushTime,
		Help: lastPushTimeHelp,
	})
	bellhopDuration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: bellhopLatency,
		Help: bellhopLatencyHelp,
	})
}

func SetPushMetrics(quoteTime int64) string {
	localMilli := getQuoteTimezoneMilli()

	if quoteTime > localMilli {
		return "Quote time more than push time"
	}

	pushDiff := localMilli - quoteTime
	lastPushTimeDuration.Set(float64(pushDiff))

	if pushDiff > settings.LongPushTime {
		longPushAlertCount.Inc()
	}

	return ""
}

func SetBellhopMetrics(sendTime int64) {
	nowMilli := time.Now().UnixMilli()
	bellhopDiff := nowMilli - sendTime

	bellhopDuration.Set(float64(bellhopDiff))
}

func getQuoteTimezoneMilli() int64 {
	// +2 часа для часового пояса котировок из редиса 	//TODO уточнить точную разницу, расхождение в несколько минут
	return time.Now().Add(time.Hour * 2).UnixMilli()
}
