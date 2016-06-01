package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/satori/go.uuid"

	log "github.com/Sirupsen/logrus"
)

var hostname string
var callCounter int64

func init() {
	hostname, _ = os.Hostname()
	callCounter = 0
}

func processStages(db *sql.DB, stdIn chan string) {
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

			go func() {

			}()

			var wg sync.WaitGroup
			if stage.Concurrency == 0 {
				stage.Concurrency = 1
			}

			wg.Add(1)
			go interrupter(&wg, &stopFlag, stdIn)

			log.Printf("Concurrency: %d", stage.Concurrency)
			wg.Add(stage.Concurrency)
			for ri := 0; ri < stage.Concurrency; ri++ {
				go worker(&wg, &stopFlag, db, &stage)
			}
			wg.Wait()
		}
		log.Printf("Stage finished!")
		if stage.Pause {
			for {
				fmt.Print(`Enter "y\\n" to continue `)
				contStr := <-stdIn
				if contStr[:1] == "y" {
					break
				}
			}
		}
	}
}

func interrupter(wg *sync.WaitGroup, stopFlag *int32, stdIn chan string) {
	fmt.Println(`Enter "n\\n" to interrupt`)
	for atomic.LoadInt32(stopFlag) > 0 {
		select {
		case str := <-stdIn:
			if str[:1] == "n" {
				atomic.StoreInt32(stopFlag, 0)
				break
			}
		}
		time.Sleep(time.Millisecond)
	}
	wg.Done()
}

func processRunOnceQueries(db *sql.DB, stage *Stage, data *QueryData) {
	for _, query := range stage.RunOnce {
		err := callTheQuery(db, query, data, query.Params)
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
		numTries := int64(1)
		if query.RandRepeat > 0 {
			numTries = rand.Int63n(int64(query.RandRepeat))
			if numTries == 0 {
				numTries++
			}
		}
		for i := int64(0); i < numTries; i++ {
			err := callTheQuery(db, query, data, query.Params)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func generateParam(paramDescriptor *Param) interface{} {
	switch paramDescriptor.Type {
	case "string":
		switch paramDescriptor.Generator {
		case "RandUUID":
			return uuid.NewV4().String()
		case "Rand4KText":
			return "AAAAAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaaAAAAaaaaaaaaaaaaaaAAAAAAAAAAAAAAAAAAAAAAAAAaaaaaaaaaa"
		default:
			panic(fmt.Sprintf("Unknown generator specified: %s", paramDescriptor.Generator))
		}
	case "int":
		switch paramDescriptor.Generator {
		case "RandToFirstQueryCallCounter":
			val := rand.Int63n(callCounter)
			if val == 0 {
				val++
			}
			return val
		default:
			panic(fmt.Sprintf("Unknown generator specified: %s", paramDescriptor.Generator))
		}
	case "timestamp":
		switch paramDescriptor.Generator {
		case "Now":
			return time.Now()
		default:
			panic(fmt.Sprintf("Unknown generator specified: %s", paramDescriptor.Generator))
		}
	default:
		panic(fmt.Sprintf("Unknown parameter type specified: %s", paramDescriptor.Type))
	}
}

func callTheQuery(db *sql.DB, query *Query, data *QueryData, query_params []*Param) error {
	params := make([]interface{}, 0, len(query_params))
	for _, query_param := range query_params {
		//		params = append(params, getFieldByName(data, query_param))
		params = append(params, generateParam(query_param))
	}

	if query.Update {
		_, err := db.Exec(query.SQL, params...)
		if err != nil {
			panic(err)
		}
	} else {
		rows, err := db.Query(query.SQL, params...)
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
	NumSQLQueries.WithLabelValues(hostname, query.QueryName).Inc()
	// TODO
	if query.QueryName == "simpleInsertUser" {
		atomic.AddInt64(&callCounter, 1)
	}

	return nil
}
