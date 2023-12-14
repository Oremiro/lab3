package main

import (
	"lab3/internal/api"
	"net/http"
)

func main() {
	startup.Run()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
