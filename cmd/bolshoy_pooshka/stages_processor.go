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
		if stage.Duration != 0 {
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

func processRunOnceQueries(stage *Stage) {
	for i, query := range stage.RunOnce {
		log.Printf("Executing a run-once query #%d: %s", i, query.SQL)
	}
}

func worker(wg *sync.WaitGroup, stage *Stage) {
	stopFlag := 0
	go func() {
		time.Sleep(stage.Duration)
		log.Printf("Setting the stopflag")
		stopFlag = 1 // No locking because this operation should be practically atomic enough
	}()
	for {
		runSingleRepeatableQuery(stage)
		if stopFlag != 0 {
			break
		}
	}
	wg.Done()
}

func runSingleRepeatableQuery(stage *Stage) {
	time.Sleep(100 * time.Millisecond)
}
