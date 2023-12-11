package main

import (
	controller "notesServer/controllers/stdhttp"
	"notesServer/gate/storage"

	"fmt"
)

func main() {

	storage := storage.NewMap()

	ADDR := "127.0.0.1:8000"

	fmt.Println("server started on ", ADDR)
	ctr := controller.NewController(ADDR, storage)

	ctr.Srv.ListenAndServe()
}
