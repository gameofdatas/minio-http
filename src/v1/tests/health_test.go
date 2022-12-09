package tests

import (
	"io"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/minio-rest/src/minio"
	routes "github.com/minio-rest/src/v1"
	"github.com/stretchr/testify/assert"
)

func Test_healthCheckHandler(t *testing.T) {
	tests := []struct {
		description   string
		route         string
		method        string
		body          io.Reader
		expectedError bool
		expectedCode  int
	}{
		{
			description:   "health check return positive response",
			route:         "/api/health",
			method:        "GET",
			body:          nil,
			expectedCode:  200,
			expectedError: false,
		},
	}

	app := fiber.New()
	ctrl := gomock.NewController(t)
	svc := minio.NewMockObjectStore(ctrl)
	routes.AddRoutes(app, minio.Service{MinioClient: &minio.ObjectStoreClient{Client: svc}})
	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.route, test.body)
		resp, err := app.Test(req, -1) // the -1 disables request latency
		assert.Equalf(t, test.expectedError, err != nil, test.description)
		if test.expectedError {
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
		// Check the response body is what we expect.
		expected := `"alive":true`
		if strings.Contains(string(body), expected) != true {
			t.Errorf("handler returned unexpected body: got %v want %v",
				string(body), expected)
		}
	}
}
