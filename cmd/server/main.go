package main

import (
	"fmt"

	"github.com/just-mitch/clogd/internal/log"
)

func main() {
	// srv := server.NewHTTPServer(":8080")
	// log.Fatal(srv.ListenAndServe())
	c := log.Config{}
	log, err := log.NewLog("/Users/mitch/apps/clogd/test-log", c)
	if err != nil {
		panic(err)
	}

	hash, err := log.Hash()
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)

	log.Close()

}
