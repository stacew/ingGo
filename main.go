package main

import (
	"log"
	"net/http"
	myapp "stacew/teamgoing/myApp"
	"stacew/teamgoing/port"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// log.Println(runtime.GOMAXPROCS(1))
	// log.Println(runtime.GOMAXPROCS(1))

	dbConn := "./test.db"
	appHandler := myapp.MakeNewHandler(dbConn)
	defer appHandler.Close()
	appHandler.Start()

	port := port.GetPort()
	log.Println("**************** ListenAndServe():" + port)
	log.Fatal(http.ListenAndServe(":"+port, appHandler)) //http start
}
