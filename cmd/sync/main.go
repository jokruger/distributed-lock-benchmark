package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	natsmutex "github.com/jokruger/nats-mutex"
	pgxmutex "github.com/jokruger/pgx-mutex"
)

func main() {
	typePtr := flag.String("type", "nats", "mutex type (nats or pgx)")
	delayPtr := flag.Duration("delay", 1000*time.Millisecond, "delay between locks")
	workersPtr := flag.Int("workers", 3, "number of workers")

	flag.Parse()

	switch *typePtr {
	case "nats":
		m, err := natsmutex.NewSyncMutex(
			natsmutex.WithUrl("nats://localhost:4222"),
			natsmutex.WithResourceID("test"),
			natsmutex.WithOwnershipCheck(false),
		)
		if err != nil {
			fmt.Println("Error creating mutex", err)
			os.Exit(1)
		}
		fmt.Println("Using", *typePtr)
		run(m, *delayPtr, *workersPtr)
		runtime.Goexit()

	case "pgx":
		m, err := pgxmutex.NewSyncMutex(
			pgxmutex.WithConnStr("postgres://postgres:postgres@localhost:5432/postgres"),
			pgxmutex.WithResourceID(123),
		)
		if err != nil {
			fmt.Println("Error creating mutex", err)
			os.Exit(1)
		}
		fmt.Println("Using", *typePtr)
		run(m, *delayPtr, *workersPtr)
		runtime.Goexit()

	default:
		fmt.Println("Invalid mutex type", typePtr)
	}
}

func run(m sync.Locker, delay time.Duration, workers int) {
	pid := os.Getpid()
	for i := 0; i < workers; i++ {
		go worker(m, delay, pid, i)
	}
}

func worker(m sync.Locker, delay time.Duration, pid int, worker int) {
	for {
		fmt.Println(time.Now(), pid, worker, "Locking")
		m.Lock()
		fmt.Println(time.Now(), pid, worker, "Processing")
		time.Sleep(delay)
		fmt.Println(time.Now(), pid, worker, "Unlocking")
		m.Unlock()
		fmt.Println(time.Now(), pid, worker, "Done")
	}
}
