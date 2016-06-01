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
)

func init() {
	//	hostname, _ := os.Hostname()
	prometheus.MustRegister(NumSQLQueries)

	//	NumSQLQueries.WithLabelValues(hostname, "wsHandler").Set(0)
}
