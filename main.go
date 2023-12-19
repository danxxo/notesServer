package main

import (
	"flag"
	controller "notesServer/controllers/stdhttp"
	Storage "notesServer/gate/storage"
	errorLogger "notesServer/pkg/errorLogger"

	"fmt"
)

func main() {

	var storage Storage.Storage

	listFlag := flag.Bool("list", false, "using list as storage")
	mapFlag := flag.Bool("map", false, "using map as storage")

	flag.Parse()

	if *listFlag {
		storage = Storage.NewList()
		fmt.Println("Using list as storage")
	} else if *mapFlag {
		storage = Storage.NewMap()
		fmt.Println("Using mpa as storage")
	} else {
		storage = Storage.NewMap()
		fmt.Println("Using mpa as storage")
	}

	ADDR := "127.0.0.1:8000"

	logger, err := errorLogger.NewErrorLogger(".log")
	if err != nil {
		fmt.Println(err)
		return
	}

	ctr := controller.NewController(ADDR, storage, logger)
	fmt.Println("server started on ", ADDR)

	ctr.Srv.ListenAndServe()
}
