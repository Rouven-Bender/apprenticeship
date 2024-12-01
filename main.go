package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
		case "addUser":
			addUser()
			break
		case "deactivateExpired":
			deactivateExpiredEntries()
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
	server := NewAPIServer(":3000", store)
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

func deactivateExpiredEntries() {
	store, err := NewSqliteStore()
	if err != nil {
		panic(err)
	}
	if err := store.DeactivateExpiredLicenses(); err != nil {
		log.Fatalln(err)
	}
}

func addUser() {
	fs := flag.NewFlagSet("addUser", flag.ExitOnError)
	var username = fs.String("u", "", "Username to be added")

	fs.Parse(os.Args[2:])

	fmt.Println("Enter the Password the User should have")
	reader := bufio.NewReader(os.Stdin)
	pwd, err := reader.ReadString('\n')
	pwd = strings.TrimSuffix(pwd, "\n")
	if err != nil {
		log.Panic(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if err != nil {
		log.Panic(err)
	}
	store, err := NewSqliteStore()
	if err != nil {
		log.Panic(err)
	}
	err = store.CreateLoginCredentials(*username, hash)
	if err != nil {
		log.Panic(err)
	}
}
