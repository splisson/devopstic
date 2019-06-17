package devopstic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/handlers"
	"github.com/splisson/devopstic/persistence"
	"github.com/splisson/devopstic/representations"
	"github.com/splisson/devopstic/services"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/http/httptest"
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
		"password": "admin",
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
	db.AutoMigrate(&entities.Commit{})
	db.AutoMigrate(&entities.Incident{})
	commitStore := persistence.NewCommitStoreDB(db)
	eventStore := persistence.NewEventDBStore(db)
	eventService := services.NewEventService(eventStore)
	commitService := services.NewCommitService(commitStore)
	commitHandlers := handlers.NewCommitHandlers(commitService)
	eventHandlers := handlers.NewEventHandlers(eventService, commitService)
	r = BuildEngine(commitHandlers, eventHandlers)
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
		authenticate(req)

		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			// Test that the http status code is 200
			statusOK := w.Code == http.StatusOK
			var pageOK bool
			if statusOK {
				var payload representations.CommitResults
				err := json.NewDecoder(w.Body).Decode(&payload)
				pageOK = err == nil && payload.Items != nil
			}

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
			"type":        entities.EVENT_COMMIT,
			"status":      entities.STATUS_SUCCESS,
			"commit_id":   uuid.New().String(),
			"pipeline_id": fmt.Sprintf("api_%s", uuid.New().String()),
			"environment": "dev",
			"timestamp":   time.Now().Unix(),
		}
		bytesRepresentation, _ := json.Marshal(message)
		body := bytes.NewBuffer(bytesRepresentation)
		req, _ = http.NewRequest("POST", "/events", body)
		req.Header.Set("Content-Type", "application/json")
		authenticate(req)

		testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
			// Test that the http status code is 200
			statusOK := w.Code == http.StatusOK
			var pageOK bool
			if statusOK {
				var result representations.Commit
				err := json.NewDecoder(w.Body).Decode(&result)
				pageOK = err == nil && result.Id != ""
			}

			return statusOK && pageOK
		})
	})

	t.Run("Post event committed and deploy", func(t *testing.T) {

		message := map[string]interface{}{
			"type":        entities.EVENT_COMMIT,
			"status":      entities.STATUS_SUCCESS,
			"commit_id":   uuid.New().String(),
			"pipeline_id": fmt.Sprintf("api_%s", uuid.New().String()),
			"environment": "dev",
			"timestamp":   time.Now().Add(mult * time.Minute).Unix(), //.Format(time.RFC3339),
		}
		message["commit_id"] = uuid.New().String()
		postEvent(t, "header", message)
		message["type"] = "approve"
		message["timestamp"] = time.Now().Unix() //.Format(time.RFC3339)
		event := postEvent(t, "header", message)
		assert.True(t, event.Type == entities.EVENT_APPROVE, "approve event")

	})
}

func postEvent(t *testing.T, authMethod string, message map[string]interface{}) representations.Event {

	bytesRepresentation, _ := json.Marshal(message)
	body := bytes.NewBuffer(bytesRepresentation)

	var req *http.Request

	if authMethod == "webhook" {
		req, _ = http.NewRequest("POST", fmt.Sprintf("/webhook/%s/events", token), body)
	} else {
		req, _ = http.NewRequest("POST", "/events", body)
		authenticate(req)
	}
	req.Header.Set("Content-Type", "application/json")
	var result representations.Event

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		err := json.NewDecoder(w.Body).Decode(&result)
		pageOK := err == nil && result.Id != ""

		return statusOK && pageOK
	})

	return result
}
