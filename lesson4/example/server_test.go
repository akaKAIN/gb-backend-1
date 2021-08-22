package example

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHandlerRoot(t *testing.T) {
	name := "Ivan"
	query := fmt.Sprintf("?name=%s", name)
	req, err := http.NewRequest("GET", "/"+query, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := &HandlerRoot{}

	handler.ServeHTTP(rr, req)

	resultStatus := rr.Code
	expectStatus := http.StatusOK
	assert.Equal(t, expectStatus, resultStatus, "Must be equal")

	data, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, bytes.Contains(data, []byte(name)), "Must contains name value in query")
}
