package main

import (
	controller "notesServer/controllers/stdhttp"
	"notesServer/gate/storage"
	errorLogger "notesServer/pkg/errorLogger"

	"fmt"
)

func main() {

	storage := storage.NewMap()

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
