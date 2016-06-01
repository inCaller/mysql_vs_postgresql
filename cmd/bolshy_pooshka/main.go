package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
)

func main() {
	_ = mysql.Config{}     // just to satisfy a bloody goimports
	_ = pq.ErrNotSupported // just to satisfy a bloody goimports

	cmdline.Parse()
	spew.Fprintf(os.Stderr, "started: %#v\n", cmdline)

	log.Infof("Starting")
	var wg sync.WaitGroup

	content, err := ioutil.ReadFile(*cmdline.Config)
	if err != nil {
		log.Fatalf("Problem reading configuration file: %v", err)
	}
	err = yaml.Unmarshal(content, &globalConfig)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %v", err)
	}
	wg.Add(1)

	test, _ := yaml.Marshal(&globalConfig)
	log.Infof("YAML: %v", string(test))
	http.Handle("/metrics", prometheus.Handler())

	db, err := sql.Open(globalConfig.DbDriver, globalConfig.DataSource)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(128)
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	stdIn := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			str, _ := reader.ReadString('\n')
			stdIn <- str
		}
	}()

	go processStages(db, stdIn)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", "0.0.0.0", 8084), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	wg.Wait()
}
