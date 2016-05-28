package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
	"time"
)

func TestCreateScenarioNumberOne(t *testing.T) {
	log.Println("Testing scenario #1")
	configNrOne := Config{
		Stages: []Stage{
			Stage{
				StageName: "doSomething",
				RPS:       100,
				Duration:  time.Second * 20,
				RunOnce: []Request{
					Request{
						RequestName: "WhoKnows",
						Query:       "SELECT 1 FROM DUAL",
						Probability: 0,
					},
				},
			},
			Stage{
				StageName: "doSomethingElse",
				RPS:       300,
				Duration:  time.Second * 60,
				RunOnce: []Request{
					Request{
						RequestName: "WhoKnows",
						Query:       "SELECT 1 FROM DUAL",
						Probability: 0,
					},
				},
				Repeat: []Request{
					Request{
						RequestName: "doSomethingUseful",
						Query:       "SELECT 1 FROM DUAL",
						Probability: 10000,
					},
					Request{
						RequestName: "doSomethingHarmful",
						Query:       "DELETE FROM table1",
						Probability: 1,
					},
				},
			},
		},
	}
	test, _ := yaml.Marshal(&configNrOne)
	log.Infof("YAML: %v", string(test))
}

func TestReadScenarioNumberTwo(t *testing.T) {
	log.Println("Testing scenario #2")
	var cfg Config

	content, err := ioutil.ReadFile("./test_config.yml")
	if err != nil {
		log.Fatalf("Problem reading configuration file: %v", err)
	}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %v", err)
	}

	test, _ := yaml.Marshal(&cfg)
	log.Infof("YAML: %v", string(test))
}
