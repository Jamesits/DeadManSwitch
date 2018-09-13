package main

import (
	"flag"
	"github.com/alexflint/go-filemutex"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"
)

func main() {
	// Parse commandline
	confPath := flag.String("conf", "config.toml", "Configuration file")
	flag.Parse()

	conf, err := loadConfig(*confPath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Dead Man's Switch starting...")

	filename := "/var/run/dmswitch.lock"
	if runtime.GOOS == "windows" {
		filename = os.TempDir() + string(os.PathSeparator) + "dmswitch.lock"
	}

	globalMutex, err := filemutex.New(filename)
	if err != nil {
		log.Fatalln("Directory did not exist or file could not created")
	}

	log.Println("Trying to acquire global lock...")
	globalMutex.Lock()
	defer globalMutex.Unlock()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	var checkTimer *time.Timer
	var uncertainCount uint = 0
	for {
		log.Println("Start routine check...")
		ret := check(conf)

		switch ret {
		case True:
			execute(conf)
		case False:
			uncertainCount = 0
		case Uncertain:
			uncertainCount++
			log.Printf("Uncertain count %d/%d", uncertainCount, conf.MaxUncertainTolerance)
			if uncertainCount > conf.MaxUncertainTolerance {
				execute(conf)
			}
		}

		checkTimer = time.NewTimer(time.Duration(conf.CheckInterval) * time.Second)
		log.Println("Idle...")
		select {
		case <- checkTimer.C:
			continue
		case <- signalChan:
			log.Println("SIGINT received, quitting...")
			os.Exit(0)
		}
	}
}