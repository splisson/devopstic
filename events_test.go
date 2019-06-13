package opstic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/splisson/opstic/entities"
	"github.com/splisson/opstic/handlers"
	"github.com/splisson/opstic/persistence"
	"github.com/splisson/opstic/services"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	r     *gin.Engine
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

func Login() (string, error) {
	message := map[string]interface{}{
		"username": "admin",
		"password": "w3yv",
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

	err := json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	token := fmt.Sprintf("%v", result["token"])
	return token, nil
}

func authenticate(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func TestMain(m *testing.M) {
	db := persistence.NewPostgresqlConnectionLocalhost()
	db.AutoMigrate(&entities.Event{})
	eventStore := persistence.NewEventDBStore(db)
	eventService := services.NewEventService(eventStore)
	eventHandlers := handlers.NewEventHandlers(eventService)
	r = BuildEngine(eventHandlers)
	var err error
	token, err = Login()
	if err != nil {
		panic(err)
	}
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

	t.Run("Post event deploy without and with authorization", func(t *testing.T) {

		//r := BuildEngine()

		// Create a request to send to the above route
		req, _ := http.NewRequest("POST", "/events", nil)

		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			// Test that the http status code is 401 because of missing authentication
			statusOK := w.Code == http.StatusUnauthorized

			return statusOK
		})

		message := map[string]interface{}{
			"category":    "deploy",
			"status":      "success",
			"commit":      "123456",
			"pipeline_id": "api",
			"environment": "dev",
			"timestamp":   time.Now().Format(time.RFC3339),
		}
		bytesRepresentation, _ := json.Marshal(message)
		body := bytes.NewBuffer(bytesRepresentation)
		req, _ = http.NewRequest("POST", "/events", body)
		req.Header.Set("Content-Type", "application/json")
		authenticate(req)

		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			// Test that the http status code is 200
			statusOK := w.Code == http.StatusOK

			//p, err := ioutil.ReadAll(w.Body)
			var result entities.Event //map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&result)
			pageOK := err == nil && result.CreatedAt.String() != ""

			return statusOK && pageOK
		})
	})

	t.Run("Post events", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())
		random := rand.Intn(10)
		mult := time.Duration(-5 * random)
		message := map[string]interface{}{
			"category":    "build",
			"status":      "success",
			"commit":      "123456",
			"pipeline_id": "api",
			"environment": "dev",
			"timestamp":   time.Now().Add(mult * time.Minute).Format(time.RFC3339),
		}
		message["commit"] = uuid.New().String()
		postEvent(t, message)
		message["category"] = "deploy"
		message["timestamp"] = time.Now().Format(time.RFC3339)
		postEvent(t, message)
	})

	t.Run("Post event via webhook", func(t *testing.T) {

		//r := BuildEngine()

		// Create a request to send to the above route

		message := map[string]interface{}{
			"category":    "incident",
			"status":      "success",
			"commit":      "123456910",
			"pipeline_id": "api",
			"environment": "dev",
			"timestamp":   time.Now().Format(time.RFC3339),
		}
		bytesRepresentation, _ := json.Marshal(message)
		body := bytes.NewBuffer(bytesRepresentation)
		req, _ := http.NewRequest("POST", fmt.Sprintf("/webhook/%s/events", token), body)
		req.Header.Set("Content-Type", "application/json")

		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			// Test that the http status code is 200
			statusOK := w.Code == http.StatusOK

			//p, err := ioutil.ReadAll(w.Body)
			var result entities.Event //map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&result)
			pageOK := err == nil && result.CreatedAt.String() != ""

			return statusOK && pageOK
		})
	})
}

func postEvent(t *testing.T, message map[string]interface{}) {

	bytesRepresentation, _ := json.Marshal(message)
	body := bytes.NewBuffer(bytesRepresentation)
	req, _ := http.NewRequest("POST", "/events", body)
	req.Header.Set("Content-Type", "application/json")
	authenticate(req)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		//p, err := ioutil.ReadAll(w.Body)
		var result entities.Event //map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&result)
		pageOK := err == nil && result.CreatedAt.String() != ""

		return statusOK && pageOK
	})
}
