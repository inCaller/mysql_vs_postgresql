package main

import (
	"bytes"
	"math/rand"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

func processStages() {
	for i, stage := range globalConfig.Stages {
		log.Printf("Started processing stage #%d, %s", i, stage.StageName)
		data := &QueryData{}
		processRunOnceQueries(&stage, data.Init())

		totalProb := float32(0)
		for _, query := range stage.Repeat {
			totalProb += query.Probability
			query.Probability = totalProb
		}
		for _, query := range stage.Repeat {
			query.Probability = query.Probability / totalProb
		}

		if stage.Duration != 0 && len(stage.Repeat) > 0 {
			var wg sync.WaitGroup
			if stage.Concurrency == 0 {
				stage.Concurrency = 1
			}
			log.Printf("Concurrency: %d", stage.Concurrency)
			wg.Add(stage.Concurrency)
			for ri := 0; ri < stage.Concurrency; ri++ {
				go worker(&wg, &stage)
			}
			wg.Wait()
		}
		log.Printf("Stage finished!")
	}
}

func processRunOnceQueries(stage *Stage, data interface{}) {
	var buf bytes.Buffer
	for i, query := range stage.RunOnce {
		buf.Reset()
		err := query.SQL.Execute(&buf, data)
		if err != nil {
			panic(err)
		}
		log.Printf("Executing a run-once query #%d: %s", i, buf.String())
	}
}

func worker(wg *sync.WaitGroup, stage *Stage) {
	stopFlag := false
	watchdog := time.AfterFunc(
		stage.Duration,
		func() {
			log.Printf("Setting the stopflag")
			stopFlag = true
		},
	)
	_ = watchdog
	for !stopFlag {
		data := &QueryData{}
		runSingleRepeatableQuery(stage, data.Init())
	}
	wg.Done()
}

func runSingleRepeatableQuery(stage *Stage, data interface{}) {
	probability := rand.Float32()
	var buf bytes.Buffer

	for i, query := range stage.Repeat {
		log.Printf("Examining a repeatable query #%d: %f %f %s", i, probability, query.Probability, query.SQL.Name())
		if query.Probability > probability {
			err := query.SQL.Execute(&buf, data)
			if err != nil {
				panic(err)
			}
			log.Printf("Executing a repeatable query #%d: %s", i, buf.String())
			return
		}
	}
	time.Sleep(100 * time.Millisecond)
}
