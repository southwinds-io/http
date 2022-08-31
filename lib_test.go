/*
  http server library - © 2018-Present SouthWinds Tech Ltd - www.southwinds.io
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package http

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestBasicToken(t *testing.T) {
	token := BasicToken("ab.cd", "1234")
	uname, pwd := ReadBasicToken(token)
	fmt.Println(uname, pwd)
}

func TestTokenOverHttp(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		t.Fatal(err)
	}
	token := BasicToken("test.user", "test-pwd")
	// add authorization token to http request headers
	req.Header.Add("Authorization", token)
	// create http client with timeout
	client := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}
	// issue http request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.StatusCode > 299 {
		t.Fatal(fmt.Errorf("response error code: %d", resp.StatusCode))
	}
}

func TestRndPwd(t *testing.T) {
	fmt.Println(RandomPwd(50, true))
}
