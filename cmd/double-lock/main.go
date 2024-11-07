package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	natsmutex "github.com/jokruger/nats-mutex"
	pgxmutex "github.com/jokruger/pgx-mutex"
)

func main() {
	typePtr := flag.String("type", "nats", "mutex type (nats or pgx)")

	flag.Parse()

	switch *typePtr {
	case "nats":
		m, err := natsmutex.NewSyncMutex(
			natsmutex.WithUrl("nats://localhost:4222"),
			natsmutex.WithResourceID("test"),
		)
		if err != nil {
			fmt.Println("Error creating mutex", err)
			os.Exit(1)
		}
		fmt.Println("Using", *typePtr)
		run(m)

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
		run(m)

	default:
		fmt.Println("Invalid mutex type", typePtr)
	}
}

func run(m sync.Locker) {
	fmt.Println("Locking first time")
	m.Lock()

	fmt.Println("Locking second time")
	m.Lock()

	fmt.Println("Done")
}
