package main

import (
	"flag"
	"log"
	gohttp "net/http"

	"github.com/gorilla/mux"

	"github.com/alka/supermart/http"
)

var port = flag.String("port", "8081", "port to listen")

func main() {
	flag.Parse()
	router := mux.NewRouter()

	http.InstallRoutes(router)

	log.Println("starting http server, listening on port:", *port)
	if err := gohttp.ListenAndServe(":"+*port, router); err != nil {
		log.Fatalf("error in starting server: %v", err)
	}

}
