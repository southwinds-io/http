/*
  http server library - © 2018-Present SouthWinds Tech Ltd - www.southwinds.io
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"testing"
	"time"
)

// launch the server
func TestServer_Serve(t *testing.T) {
	// create a new server
	s := New("test", "1.0")
	// set auth credentials
	err := os.Setenv("HTTP_USER", "test.user")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("HTTP_PASSWORD", "test-pwd")
	if err != nil {
		t.Fatal(err)
	}
	// configure http handlers
	s.Http = func(router *mux.Router) {
		// add middlewares
		router.Use(s.AuthenticationMiddleware)
		router.Use(s.LoggingMiddleware)
		// add handler
		router.HandleFunc("/", doSomething).Methods("GET")
		router.HandleFunc("/items/abc", doSomething).Methods("GET")
	}
	// customised authentication on a path basis
	s.Auth = map[string]func(http.Request) *UserPrincipal{
		// authorise all requests for the path /items/....
		"^/items/.*": func(request http.Request) *UserPrincipal {
			return &UserPrincipal{
				Username: "guest",
				Created:  time.Now(),
			}
		},
	}
	// serve
	s.Serve()
}

func doSomething(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
