package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"testing"
	"time"
)

func TestScenarioNumberOne(t *testing.T) {
	log.Println("Testing scenario #1")
	configNrOne := Config{
		Stages: []Stage{
			Stage{
				StageName: "doSomething",
				RPS:       100,
				Duration:  time.Second * 20,
				RunOnce: []Request{
					Request{
						Name:        "WhoKnows",
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
						Name:        "WhoKnows",
						Query:       "SELECT 1 FROM DUAL",
						Probability: 0,
					},
				},
				Repeat: []Request{
					Request{
						Name:        "doSomethingUseful",
						Query:       "SELECT 1 FROM DUAL",
						Probability: 10000,
					},
					Request{
						Name:        "doSomethingHarmful",
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
