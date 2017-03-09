package main

import (
	"fmt"
	"log"
	"net/http"
)

func init() {
	http.Handle("/", middleWare(http.HandlerFunc(index)))
}

func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//code before
		log.Println("Middleware mulai")
		next.ServeHTTP(w, r)
		//code after
		log.Println("Middleware end")
	})
}
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<p>success</p>")
}
