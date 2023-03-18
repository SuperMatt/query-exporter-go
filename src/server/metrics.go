package server

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func MerticName(name string) string {
	return fmt.Sprintf("%s_%s", prefix, name)
}

var (
	prefix = "query_exporter"
	// Metrics

	// Info about the queries
	queryInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: MerticName("query_info"),
		Help: "Information about the query",
	}, []string{"address", "name", "query", "query_offset"})

	// timestamp of the query
	queryTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: MerticName("query_timestamp"),
		Help: "The timestamp of the query",
	}, []string{"name", "query_offset"})

	// the query status
	queryStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: MerticName("query_status"),
		Help: "The status of the query",
	}, []string{"name", "query_offset"})

	// the query duration
	queryDuration = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: MerticName("query_duration"),
		Help: "The duration of the query",
	}, []string{"name", "query_offset"})

	// if the query was successful
	querySuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: MerticName("query_success"),
		Help: "If the query was successful",
	}, []string{"name", "query_offset"})

	queryDnsError = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: MerticName("query_dns_error"),
		Help: "If the query had a DNS error",
	}, []string{"name"})

	queryPortError = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: MerticName("query_port_error"),
		Help: "If the query had a port error",
	}, []string{"name"})
)
