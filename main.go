package main

import (
	"flag"
	"log"
	"loveoneanother.at/tiedot/db"
	"loveoneanother.at/tiedot/srv/v1"
	"loveoneanother.at/tiedot/srv/v2"
	"os"
	"runtime"
	"strconv"
)

func main() {
	var defaultMaxprocs int
	var err error
	if defaultMaxprocs, err = strconv.Atoi(os.Getenv("GOMAXPROCS")); err != nil {
		defaultMaxprocs = runtime.NumCPU() * 2
	}
	var mode, dir string
	var port, maxprocs int
	flag.StringVar(&mode, "mode", "", "[v1|v2|bench|durable-bench|example]")
	flag.StringVar(&dir, "dir", "", "database directory")
	flag.IntVar(&port, "port", 0, "listening port number")
	flag.IntVar(&maxprocs, "gomaxprocs", defaultMaxprocs, "GOMAXPROCS")
	flag.Parse()

	if mode == "" {
		flag.PrintDefaults()
		return
	}

	runtime.GOMAXPROCS(maxprocs)
	log.Printf("GOMAXPROCS is set to %d", maxprocs)

	if maxprocs < runtime.NumCPU() {
		log.Printf("GOMAXPROCS (%d) is less than number of CPUs (%d), this may affect performance. You can change it via environment variable GOMAXPROCS or by passing CLI parameter -gomaxprocs", maxprocs, runtime.NumCPU())
	}

	switch mode {
	case "v1":
		fallthrough
	case "v2":
		if dir == "" {
			log.Fatal("Please specify database directory, for example -dir=/tmp/db")
		}
		if port == 0 {
			log.Fatal("Please specify port number, for example -port=8080")
		}
		db, err := db.OpenDB(dir)
		if err != nil {
			log.Fatal(err)
		}
		if mode == "v1" {
			v1.Start(db, port)
		} else if mode == "v2" {
			v2.Start(db, port)
		}
	case "bench":
		benchmark()
	case "durable-bench":
		durableBenchmark()
	case "example":
		embeddedExample()
	default:
		flag.PrintDefaults()
		return
	}
}
