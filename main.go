package main

import (
	"log"
	"net/http"
	"proxy/pys"
)

func main() {
	http.HandleFunc("/pys-pixel-proxy", pys.Handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:9091", nil))
}
