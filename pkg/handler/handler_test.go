package handler

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/nais/armorator/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetPolicy(t *testing.T) {
	var app = NewApp(context.Background(), &config.Config{
		DevelopmentMode: true,
		Port:            "8090",
	}, logrus.NewEntry(logrus.New()))
	app.SetupHttpRouter(NewHandler(app))

	req, err := http.NewRequest(http.MethodGet, EndpointGetPolicy, nil)
	assert.NoError(t, err)

	response := executeRequest(req, app.Router)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request, router *mux.Router) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
