package opstic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	r *gin.Engine
	token string
)

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}


func Login() string {
	message := map[string]interface{}{
		"username": "admin",
		"password":  "w3yv",
	}
	bytesRepresentation, _ := json.Marshal(message)
	body := bytes.NewBuffer(bytesRepresentation)

	req, _ := http.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")
	var result map[string]interface{}

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	json.NewDecoder(w.Body).Decode(&result)

	token := fmt.Sprintf("%v", result["token"])
	return token
}

func authenticate(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func TestMain(m *testing.M) {
	r = BuildEngine()
	token = Login()
	m.Run()
}

func TestGetEvents(t *testing.T) {

	t.Run("should get events in the response", func(t *testing.T) {

		// Create a request to send to the above route
		req, _ := http.NewRequest("GET", "/events", nil)

		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			// Test that the http status code is 200
			statusOK := w.Code == http.StatusOK

			// Test that the page title is "Home Page"
			// You can carry out a lot more detailed tests using libraries that can
			// parse and process HTML pages
			p, err := ioutil.ReadAll(w.Body)
			pageOK := err == nil && strings.Index(string(p), "events") > 0

			return statusOK && pageOK
		})
	})
}

func TestPostEvents(t *testing.T) {

	t.Run("Post events", func(t *testing.T) {

		//r := BuildEngine()

		// Create a request to send to the above route
		req, _ := http.NewRequest("POST", "/events", nil)

		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			// Test that the http status code is 401 because of missing authentication
			statusOK := w.Code == http.StatusUnauthorized

			return statusOK
		})

		req, _ = http.NewRequest("POST", "/events", nil)
		req.Header.Set("Content-Type", "application/json")
		authenticate(req)
		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			// Test that the http status code is 200
			statusOK := w.Code == http.StatusOK

			p, err := ioutil.ReadAll(w.Body)
			pageOK := err == nil && strings.Index(string(p), "created") > 0

			return statusOK && pageOK
		})
	})
}

