package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"
)

func processStages(db *sql.DB) {
	for i, stage := range globalConfig.Stages {
		log.Printf("Started processing stage #%d, %s", i, stage.StageName)
		data := &QueryData{}
		processRunOnceQueries(db, &stage, data.Init())

		totalProb := float32(0)
		for _, query := range stage.Repeat {
			totalProb += query.Probability
			query.Probability = totalProb
		}
		for _, query := range stage.Repeat {
			query.Probability = query.Probability / totalProb
		}

		if stage.Duration != 0 && len(stage.Repeat) > 0 {
			stopFlag := int32(1)
			watchdog := time.AfterFunc(
				stage.Duration,
				func() {
					log.Printf("Setting the stopflag")
					atomic.StoreInt32(&stopFlag, 0)
				},
			)
			_ = watchdog

			var wg sync.WaitGroup
			if stage.Concurrency == 0 {
				stage.Concurrency = 1
			}
			log.Printf("Concurrency: %d", stage.Concurrency)
			wg.Add(stage.Concurrency)
			for ri := 0; ri < stage.Concurrency; ri++ {
				go worker(&wg, &stopFlag, db, &stage)
			}
			wg.Wait()
		}
		log.Printf("Stage finished!")
		if stage.Pause {
			contStr := " "
			reader := bufio.NewReader(os.Stdin)
			for contStr[:1] != "y" {
				fmt.Print("Continue? ")
				contStr, _ = reader.ReadString('\n')
			}
		}
	}
}

func processRunOnceQueries(db *sql.DB, stage *Stage, data *QueryData) {
	for _, query := range stage.RunOnce {
		err := callTheQuery(db, query.Update, query.SQL, data, query.Params)
		if err != nil {
			panic(err)
		}
	}
}

func worker(wg *sync.WaitGroup, stopFlag *int32, db *sql.DB, stage *Stage) {
	for atomic.LoadInt32(stopFlag) > 0 {
		data := &QueryData{}
		runSingleRepeatableScenario(db, stage, data.Init())
	}
	wg.Done()
}

func runSingleRepeatableScenario(db *sql.DB, stage *Stage, data *QueryData) {
	probability := rand.Float32()

	for _, scenario := range stage.Repeat {
		if scenario.Probability > probability {
			err := runScenario(db, scenario, data)
			if err != nil {
				panic(err)
			}
			return
		}
	}
}

func runScenario(db *sql.DB, scenario *Scenario, data *QueryData) error {
	for _, query := range scenario.Queries {
		err := callTheQuery(db, query.Update, query.SQL, data, query.Params)
		if err != nil {
			return err
		}
	}

	return nil
}

func callTheQuery(db *sql.DB, update bool, query string, data *QueryData, paramsNames []string) error {
	params := make([]interface{}, 0, len(paramsNames))
	for _, paramName := range paramsNames {
		params = append(params, getFieldByName(data, paramName))
	}
	log.Printf("Executing a repeatable query: %s (%q) (%#+v)", query, paramsNames, params)

	if update {
		_, err := db.Exec(query, params...)
		if err != nil {
			panic(err)
		}
	} else {
		rows, err := db.Query(query, params...)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		columnNames, err := rows.Columns()
		if err != nil {
			panic(err)
		}
		rc := NewMapStringScan(columnNames)
		for rows.Next() {
			err := rc.Update(rows)
			if err != nil {
				panic(err)
			}
			_ = rc.Get()
		}
	}

	return nil
}
