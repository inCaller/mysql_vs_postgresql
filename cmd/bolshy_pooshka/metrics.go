package main

import (
	"github.com/prometheus/client_golang/prometheus"
	//	"os"
)

var (
	NumSQLQueries = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "pooshka_sql_queries_total",
		Help: "Number of total SQL queries.",
	}, []string{"server", "query_name"})
	SQLQueriesTimes = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "pooshka_sql_queries_times_usecs",
		Help:       "SQL timers (in microseconds).",
		AgeBuckets: 3,
	}, []string{"server", "query_name"})
	SelectQueriesTimes = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "pooshka_select_queries_times_usecs",
		Help:       "Select (exec + resultset) timers (in microseconds).",
		AgeBuckets: 3,
	}, []string{"server", "query_name"})
)

func init() {
	//	hostname, _ := os.Hostname()
	prometheus.MustRegister(NumSQLQueries)
	prometheus.MustRegister(SQLQueriesTimes)
	prometheus.MustRegister(SelectQueriesTimes)

	//	NumSQLQueries.WithLabelValues(hostname, "wsHandler").Set(0)
}
