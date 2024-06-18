package main

import (
	"log"
	"os"
	"path/filepath"
	"taskm/cmd"
	"taskm/db"

	"github.com/mitchellh/go-homedir"
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	dbPath := filepath.Join(home, "tasks.db")

	must(db.Init(dbPath))

	must(cmd.Execute())
}

func must(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
