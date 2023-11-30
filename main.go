package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/izaakdale/gh-actions-go/router"
)

func main() {
	mux := router.New()
	http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")), mux)
}
