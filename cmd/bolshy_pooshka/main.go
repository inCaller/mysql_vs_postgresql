package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"sync"
)

func main() {
	log.Infof("Starting")
	var wg sync.WaitGroup

	content, err := ioutil.ReadFile("./test_config.yml")
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

	go processStages()

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", "0.0.0.0", 8084), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	wg.Wait()
}
