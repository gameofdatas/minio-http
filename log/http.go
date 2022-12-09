package log

import (
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LoggingServer struct {
	client HTTPClient
}

func (l *LoggingServer) Do(req *http.Request) (*http.Response, error) {
	return l.client.Do(req)
}

func NewLoggerClient() *LoggingServer {
	return &LoggingServer{
		client: &http.Client{},
	}
}
