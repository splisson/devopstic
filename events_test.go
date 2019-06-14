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
	db.AutoMigrate(&entities.Deployment{})
	commitStore := persistence.NewCommitStoreDB(db)
	deploymentStore := persistence.NewDeploymentDBStore(db)
	commitService := services.NewCommitService(commitStore, deploymentStore)
	commitHandlers := handlers.NewCommitHandlers(commitService)
	r = BuildEngine(commitHandlers)
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
			"type":        entities.COMMIT_EVENT_COMMIT,
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

	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(10)
	mult := time.Duration(-5 * random)

	t.Run("Post event committed and deploy", func(t *testing.T) {

		message := map[string]interface{}{
			"type":        entities.COMMIT_EVENT_COMMIT,
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
		commit := postEvent(t, "header", message)
		assert.True(t, commit.ApprovalTime > 0, "approval time > 0")
		message["commit_id"] = uuid.New().String()
		message["type"] = "commit"
		message["timestamp"] = time.Now().Add(mult * time.Minute).Unix()
		postEvent(t, "header", message)
		message["type"] = "deploy"
		message["timestamp"] = time.Now().Unix() //.Format(time.RFC3339)
		commit = postEvent(t, "header", message)
		assert.True(t, commit.TotalLeadTime > 0, "lead time > 0")
	})

	t.Run("Post event via webhook", func(t *testing.T) {

		//r := BuildEngine()

		// Create a request to send to the above route

		message := map[string]interface{}{
			"category":    "incident",
			"status":      "failure",
			"commit_id":   "123456910",
			"pipeline_id": "api",
			"environment": "dev",
			"timestamp":   time.Now().Add(mult * time.Minute).Unix(),
		}

		postEvent(t, "webhook", message)
		message["status"] = "success"
		message["timestamp"] = time.Now().Unix()
		postEvent(t, "webhook", message)

		//bytesRepresentation, _ := json.Marshal(message)
		//body := bytes.NewBuffer(bytesRepresentation)
		//req, _ := http.NewRequest("POST", fmt.Sprintf("/webhook/%s/events", token), body)
		//req.Header.Set("Content-Type", "application/json")
		//
		//testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		//	// Test that the http status code is 200
		//	statusOK := w.Code == http.StatusOK
		//
		//	//p, err := ioutil.ReadAll(w.Body)
		//	var result entities.Deployment //map[string]interface{}
		//	err := json.NewDecoder(w.Body).Decode(&result)
		//	pageOK := err == nil && result.CreatedAt.String() != ""
		//
		//	return statusOK && pageOK
		//})
	})
}

func postEvent(t *testing.T, authMethod string, message map[string]interface{}) representations.Commit {

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
	var result representations.Commit

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		// Test that the http status code is 200
		statusOK := w.Code == http.StatusOK

		//p, err := ioutil.ReadAll(w.Body)
		//map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&result)
		pageOK := err == nil && result.Id != ""

		return statusOK && pageOK
	})

	return result
}
