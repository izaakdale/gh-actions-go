package router

import (
	"context"
	"net/http"
)

const (
	CtxKey ctxKV = "testkey"
	CtxVal ctxKV = "testval"
)

type ctxKV string

func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", GetTest)
	return mux
}

func TestMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), CtxKey, CtxVal))
		next.ServeHTTP(w, r)
	})
}

func GetTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Error(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("hello"))
}

func Error(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
