package controller

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleHome(t *testing.T) {
	h := new(home)
	expectedBody := "home template"
	expectedHeader := "text/html"
	expectedStatusCode := http.StatusOK
	h.homeTemplate, _ = template.New("").Parse(expectedBody)

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Failed to create a http request. %s", err.Error())
	}
	w := httptest.NewRecorder()

	h.handleHome(w, r)

	actual, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatalf("Failed to read the body of the response. %s", err.Error())
	}

	if string(actual) != expectedBody {
		t.Errorf("Failed to execute correct template")
	}
	sc := w.Result().StatusCode
	if sc != expectedStatusCode {
		t.Errorf("Expected status: %d, got %d", expectedStatusCode, sc)
	}

	ct := w.Result().Header.Get("Content-Type")
	if !strings.Contains(ct, expectedHeader) {
		t.Errorf("Expected Content-Type Header: %s, got %s", expectedHeader, ct)
	}
}
