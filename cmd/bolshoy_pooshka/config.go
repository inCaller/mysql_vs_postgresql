package main

import "time"

type Config struct {
	Stages []Stage
}

type Stage struct {
	Name     string
	RPS      float32       // 0 - infinity
	Duration time.Duration /*
		0 - end as soon as all the RunOnce queries done
		duration - obvious
		set a huge duration to run until interrupted
	*/
	RunOnce []Request // executed one by one
	Repeat  []Request // executed in parallel according to their probability
	Pause   bool      // Do not step to the next stage automatically
}

type Request struct {
	Name        string  // used as a part of metric name
	Query       string  // SQL itself
	Probability float32 // 0 - never, 1 - each time, ignored for RunOnce
}
