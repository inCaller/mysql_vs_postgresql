package main

import (
	log "github.com/Sirupsen/logrus"
	"sync"
	"time"
)

func processStages() {
	for i, stage := range globalConfig.Stages {
		log.Printf("Started processing stage #%d, %s", i, stage.StageName)
		processRunOnceQueries(&stage)
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
		log.Printf("Stage finished!")
	}
}

func processRunOnceQueries(stage *Stage) {
	for i, query := range stage.RunOnce {
		log.Printf("Executing a run-once query #%d: %s", i, query.SQL)
	}
}

func worker(wg *sync.WaitGroup, stage *Stage) {
	time.Sleep(stage.Duration)
	wg.Done()
}
