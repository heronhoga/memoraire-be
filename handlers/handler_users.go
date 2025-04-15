package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}