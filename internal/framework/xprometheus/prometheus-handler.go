package xprometheus

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type PrometheusHandler struct {
	counter   *prometheus.CounterVec
	histogram *prometheus.HistogramVec
}

func NewPrometheusHandler(namespace *string) *PrometheusHandler {

	counterName := "api_request_number_total"
	histogramName := "api_request_duration_seconds"

	if namespace != nil {
		counterName = fmt.Sprintf("%s_%s", *namespace, counterName)
		histogramName = fmt.Sprintf("%s_%s", *namespace, histogramName)
	}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: counterName,
		Help: "Total number API requests sent by the service.",
	}, []string{"handler", "method", "status"})

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: histogramName,
		Help: "Duration of the API requests being handled",
	}, []string{"handler"})

	if err := prometheus.Register(counter); err != nil {
		logrus.WithError(err).Warn("failed to register prometheus counter")
	}

	if err := prometheus.Register(histogram); err != nil {
		logrus.WithError(err).Warn("failed to register prometheus histogram")
	}

	h := &PrometheusHandler{
		counter:   counter,
		histogram: histogram,
	}

	return h
}

func (h *PrometheusHandler) InitMetric(path, method string) {

	h.counter.WithLabelValues(path, method, strconv.Itoa(200))
	h.histogram.WithLabelValues(path)
}

func (h *PrometheusHandler) HandleRequest(path, method string, status int, duration time.Duration) {

	h.counter.WithLabelValues(path, method, strconv.Itoa(status)).Inc()
	h.histogram.WithLabelValues(path).Observe(duration.Seconds())
}
