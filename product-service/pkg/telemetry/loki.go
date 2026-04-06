package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LokiWriter implements io.Writer — each write is pushed as a log stream to Loki.
// It is used as a fan-out writer alongside stdout in the logger.
type LokiWriter struct {
	url         string
	serviceName string
	client      *http.Client
}

// NewLokiWriter creates a non-blocking Loki push writer.
func NewLokiWriter(lokiURL, serviceName string) *LokiWriter {
	return &LokiWriter{
		url:         lokiURL + "/loki/api/v1/push",
		serviceName: serviceName,
		client:      &http.Client{Timeout: 3 * time.Second},
	}
}

// Write pushes p as a single log line to the Loki push API, asynchronously.
// Fire-and-forget: errors are silently discarded to never block the application.
func (lw *LokiWriter) Write(p []byte) (n int, err error) {
	go lw.push(string(p))
	return len(p), nil
}

type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

func (lw *LokiWriter) push(line string) {
	ts := fmt.Sprintf("%d", time.Now().UnixNano())

	payload := lokiPushRequest{
		Streams: []lokiStream{
			{
				Stream: map[string]string{
					"service": lw.serviceName,
					"app":     "oolio-food-ordering-backend",
				},
				Values: [][]string{{ts, line}},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, lw.url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := lw.client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}
