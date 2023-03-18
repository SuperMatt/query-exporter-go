package server

// create a web server which will serve the metrics from the query sent to the endpoint
// the server should be able to be configured to listen on a specific port

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/supermatt/query-exporter-go/src/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Logger struct {
	info  *log.Logger
	debug *log.Logger
	error *log.Logger

	debugEnabled bool
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.debugEnabled {
		l.debug.Printf(format, v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.error.Printf(format, v...)
}

// Server is the main server object
type Server struct {
	*config.Config

	Logger Logger
}

// NewServer creates a new server object
func NewServer(cfg *config.Config) *Server {
	s := &Server{
		Config: cfg,
		Logger: Logger{
			info:         log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			debug:        log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
			error:        log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
			debugEnabled: cfg.Debug,
		},
	}
	return s
}

func (s *Server) recordMetrics() {

	go func() {
		for {

			for _, endpoint := range s.Config.Endpoints {
				if len(endpoint.QueryOffsets) == 0 {
					endpoint.QueryOffsets = []string{"Now"}
				}

				for _, queryOffset := range endpoint.QueryOffsets {

					dnsError := 0
					portErrors := 0
					statusCode := 0
					success := 0
					timestamp := time.Now().Unix()
					duration := int64(0)
					timestampParseError := 0

					queryTime, err := ParseQueryOffset(queryOffset)
					if err != nil {
						timestampParseError = 1
						s.Logger.Errorf("Error parsing query offset: %s", err)
					}

					s.Logger.Debugf("Query timestamp: %d", queryTime)

					query := "count(up)"
					if endpoint.Query != "" {
						query = endpoint.Query
					}

					// probe the port to see if it's open
					parsedUrl, err := url.Parse(endpoint.Address)
					if err != nil {
						s.Logger.Errorf("Error parsing url: %s", err)
					}

					host := parsedUrl.Hostname()
					port := parsedUrl.Port()
					scheme := parsedUrl.Scheme

					ipAddr, err := GetIP(host)

					if err != nil {
						dnsError = 1
						s.Logger.Errorf("Error getting IP address: %s", err)
					} else {
						_, err := CheckPort(ipAddr, port)
						if err != nil {
							portErrors = 1
							s.Logger.Errorf("Error checking port: %s", err)
						} else {
							success, statusCode, timestamp, duration = s.QueryPrometheus(endpoint, scheme, host, port, query, queryTime)
						}
					}
					queryDnsError.WithLabelValues(endpoint.Name).Set(float64(dnsError))
					queryPortError.WithLabelValues(endpoint.Name).Set(float64(portErrors))
					querySuccess.WithLabelValues(endpoint.Name, queryOffset).Set(float64(success))
					queryStatus.WithLabelValues(endpoint.Name, queryOffset).Set(float64(statusCode))
					queryInfo.WithLabelValues(endpoint.Address, endpoint.Name, query, queryOffset).Set(1)
					queryTimestamp.WithLabelValues(endpoint.Name, queryOffset).Set(float64(timestamp))
					queryDuration.WithLabelValues(endpoint.Name, queryOffset).Set(float64(duration))
					queryTimestampError.WithLabelValues(endpoint.Name, queryOffset).Set(float64(timestampParseError))
				}
			}
			time.Sleep(15 * time.Second)
		}
	}()
}

// Start starts the server
func (s *Server) Start() error {
	s.Logger.Infof("Starting server on port %d", s.Server.Port)
	s.Logger.Debugf("Config: %+v", s.Config)

	s.recordMetrics()

	http.Handle(s.Server.MetricsEndpoint, promhttp.Handler())
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Server.Port), nil)
}
