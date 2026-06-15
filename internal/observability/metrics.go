package observability

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// Metrics stores process counters in Prometheus text format.
type Metrics struct {
	namespace string

	mu              sync.RWMutex
	requestsTotal   map[string]uint64
	requestDuration map[string]time.Duration
	denialsTotal    map[string]uint64
}

func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		namespace:       sanitizeMetricName(namespace),
		requestsTotal:   make(map[string]uint64),
		requestDuration: make(map[string]time.Duration),
		denialsTotal:    make(map[string]uint64),
	}
}

func (m *Metrics) ObserveRequest(path string, status int, duration time.Duration) {
	key := fmt.Sprintf("path=%s,status=%d", path, status)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestsTotal[key]++
	m.requestDuration[key] += duration
}

func (m *Metrics) ObserveDenial(reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.denialsTotal[reason]++
}

func (m *Metrics) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		_, _ = w.Write([]byte(m.Render()))
	})
}

func (m *Metrics) Render() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var builder strings.Builder
	requestMetric := m.namespace + "_requests_total"
	durationMetric := m.namespace + "_request_duration_seconds_total"
	denialMetric := m.namespace + "_denials_total"

	builder.WriteString("# TYPE " + requestMetric + " counter\n")
	for _, key := range sortedKeys(m.requestsTotal) {
		builder.WriteString(fmt.Sprintf("%s{%s} %d\n", requestMetric, toLabels(key), m.requestsTotal[key]))
	}

	builder.WriteString("# TYPE " + durationMetric + " counter\n")
	for _, key := range sortedKeysDuration(m.requestDuration) {
		builder.WriteString(fmt.Sprintf("%s{%s} %.6f\n", durationMetric, toLabels(key), m.requestDuration[key].Seconds()))
	}

	builder.WriteString("# TYPE " + denialMetric + " counter\n")
	for _, key := range sortedKeys(m.denialsTotal) {
		builder.WriteString(fmt.Sprintf("%s{reason=%q} %d\n", denialMetric, key, m.denialsTotal[key]))
	}

	return builder.String()
}

func sanitizeMetricName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "aegis_mcp"
	}
	value = strings.ReplaceAll(value, "-", "_")
	return value
}

func toLabels(value string) string {
	parts := strings.Split(value, ",")
	labels := make([]string, 0, len(parts))
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		labels = append(labels, fmt.Sprintf("%s=%q", kv[0], kv[1]))
	}
	return strings.Join(labels, ",")
}

func sortedKeys(m map[string]uint64) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedKeysDuration(m map[string]time.Duration) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
