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
	ac := ActualConnection{}
	mux.HandleFunc("/test", GetTest)
	mux.HandleFunc("/users", GetUsers(ac))
	return mux
}

func SomeMiddle(next http.Handler) http.Handler {
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

type DbConn interface {
	Ping() error
}
type ActualConnection struct {
}

func (a ActualConnection) Ping() error {
	return nil
}

func GetUsers(db DbConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			Error(w, http.StatusInternalServerError)
			return
		}
		w.Write([]byte("users"))
	}
}
