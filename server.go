package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", index)

	http.ListenAndServe(":80", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "onlypage.html")
}
