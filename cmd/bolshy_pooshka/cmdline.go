package main

import (
	"flag"
	"fmt"
	"os"
)

type CmdLine struct {
	Config *string
}

var cmdline = CmdLine{
	Config: flag.String("c", "./test_config.yml", "Config file name"),
}

func (*CmdLine) Parse() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()

	}

	flag.Parse()
}

type StrList []string

func (s *StrList) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *StrList) Set(v string) error {
	*s = append(*s, v)
	return nil
}
