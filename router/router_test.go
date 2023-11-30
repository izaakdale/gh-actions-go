package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/izaakdale/gh-actions-go/router"
)

func TestMiddleware(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctxCheckFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(router.CtxKey) != router.CtxVal {
				t.Error("context was not added correctly")
			}
		})
		hand := router.TestMiddle(ctxCheckFunc)

		rr := httptest.NewRecorder()
		hand.ServeHTTP(rr, &http.Request{})
	})

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

	t.Run("failure on purpose", func(t *testing.T) {
		t.Error("testing gh action conditional when test fails")
	})
}
