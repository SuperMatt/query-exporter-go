package server

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/supermatt/query-exporter-go/src/config"
)

func (s *Server) QueryPrometheus(endpoint config.Endpoint, scheme, host, port, query string, queryTime int64) (success int, statusCode int, timestamp int64, duration int64) {
	address := url.URL{
		Scheme:   scheme,
		Host:     fmt.Sprintf("%s:%s", host, port),
		Path:     "/api/v1/query",
		RawQuery: fmt.Sprintf("query=%s&time=%d", query, queryTime),
	}
	s.Logger.Debugf("Querying: %s", address.String())

	client := &http.Client{}
	req, _ := http.NewRequest("GET", address.String(), nil)
	for _, header := range endpoint.Headers {
		req.Header.Add(header.Name, header.Value)
		s.Logger.Debugf("Adding header: %s: %s", header.Name, header.Value)
	}

	startTime := time.Now()
	resp, err := client.Do(req)
	endTime := time.Now()

	timestamp = startTime.Unix()
	duration = endTime.Sub(startTime).Milliseconds()

	if err != nil {
		s.Logger.Errorf("Error querying endpoint: %s", err)
	} else {
		success = 1
		statusCode = resp.StatusCode
	}

	return success, statusCode, timestamp, duration

}
