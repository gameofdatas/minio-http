package tests

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/minio-rest/src/minio"
	routes "github.com/minio-rest/src/v1"
	miniov7 "github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	tests := []struct {
		description        string
		route              string
		method             string
		expectedError      bool
		requestBody        io.Reader
		expectedStatusCode int
		mockService        func(*minio.MockObjectStore) *minio.MockObjectStore
		expectedBody       []byte
	}{
		{
			description:        "upload a file to minio",
			route:              "/api/upload",
			method:             "POST",
			expectedError:      false,
			requestBody:        strings.NewReader(`{"bucket":"soruce-bucket","file_path":"x-adapatation-id","file_name":"metadata.json","file_content":"c2FtcGxl","retention_ms":900000}`),
			expectedStatusCode: http.StatusCreated,
			mockService: func(store *minio.MockObjectStore) *minio.MockObjectStore {
				store.EXPECT().BucketExists(gomock.Any(), gomock.Any()).Return(true, nil)
				store.EXPECT().PutObject(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(miniov7.UploadInfo{
					Bucket: "source-bukcet",
					Key:    "x-adapatation-id/metadata.json",
				}, nil)
				return store
			},
			expectedBody: []byte(`Created`),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			app := fiber.New()
			ctrl := gomock.NewController(t)
			svc := test.mockService(minio.NewMockObjectStore(ctrl))
			routes.AddRoutes(app, minio.Service{MinioClient: &minio.ObjectStoreClient{Client: svc}})
			req := httptest.NewRequest(test.method, test.route, test.requestBody)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req, -1) // the -1 disables request latency
			assert.Equalf(t, test.expectedError, err != nil, test.description)
			if test.expectedError {
				return
			}
			// Check the status code is what we expect.
			assert.Equalf(t, test.expectedStatusCode, resp.StatusCode, test.description)

			// Check the response body is what we expect.
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fail()
			}
			assert.Equalf(t, test.expectedBody, body, test.description)
		})
	}
}
