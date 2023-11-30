package router_test

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/izaakdale/gh-actions-go/router"
)

func TestNew(t *testing.T) {
	mux := router.New()
	if mux == nil {
		t.Error("nil mux")
	}
}

func TestMiddleware(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctxCheckFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(router.CtxKey) != router.CtxVal {
				t.Error("context was not added correctly")
			}
		})
		hand := router.SomeMiddle(ctxCheckFunc)

		rr := httptest.NewRecorder()
		hand.ServeHTTP(rr, &http.Request{})
	})

}

type TestDBConn struct {
}

func (t TestDBConn) Ping() error {
	log.Printf("Hit ping\n")
	return nil
}

func TestGetter(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router.GetTest(rr, req)

		if rr.Result().StatusCode != http.StatusOK {
			t.Error("status code is incorrect")
		}

		if rr.Body.String() != "hello" {
			t.Error("response is incorrect")
		}
	})

	t.Run("sad path", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router.GetTest(rr, req)

		if rr.Result().StatusCode != http.StatusMethodNotAllowed {
			t.Error("status code is incorrect")
		}
	})

}

type PassStub struct {
}

func (p PassStub) Ping() error {
	return nil
}

type FailStub struct {
}

func (f FailStub) Ping() error {
	return errors.New("uh oh")
}

func TestUserGetter(t *testing.T) {

	t.Run("user happy", func(t *testing.T) {
		p := PassStub{}

		hand := router.GetUsers(p)

		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		hand(rr, req)

		if rr.Result().StatusCode != http.StatusOK {
			t.Error("wrong status")
		}

		if rr.Body.String() != "users" {
			t.Error("response is incorrect")
		}

	})
	t.Run("user sad", func(t *testing.T) {
		f := FailStub{}

		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		hand := router.GetUsers(f)

		hand(rr, req)

		if rr.Result().StatusCode != http.StatusInternalServerError {
			t.Error("wrong status")
		}
	})
}
