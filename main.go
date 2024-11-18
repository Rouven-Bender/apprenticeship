package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("expecting subcommand")
		os.Exit(1)
	} else {
		switch os.Args[1] {
		case "serve":
			runServer()
			break
		case "export":
			runExport()
			break
		default:
			log.Println("unknown subcommand")
		}
	}
}

func runServer() {
	store, err := NewSqliteStore()
	if err != nil {
		panic(err)
	}
	server := NewAPIServer(":3000", *store)
	server.Run()
}

func runExport() {
	store, err := NewSqliteStore()
	if err != nil {
		panic(err)
	}
	lics, err := store.GetAllActivSublicenses()
	if err != nil {
		log.Println(err)
	}
	j, err := json.Marshal(lics)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s", string(j))
}
