package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSendMetric(t *testing.T) {
	// Setup test data
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request URL and method
		if r.URL.Path != "/update/gauge/metric/1" || r.Method != "POST" {
			t.Errorf("Unexpected request URL or method: %v %v", r.URL.Path, r.Method)
		}
		// Verify request header
		contentType := r.Header.Get("Content-Type")
		if contentType != "text/plain" {
			t.Errorf("Unexpected content type: %v", contentType)
		}
		// Write a response
		fmt.Fprint(w, "OK")
	}))
	defer testServer.Close()

	type args struct {
		mType metricType
		name  string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "positive test",
			args: args{mType: "gauge", name: "metric", value: 1},
			want: nil,
		},
	}

	reportInterval := 10 * time.Second
	pollInterval := 2 * time.Second
	agent := NewAgent(reportInterval, pollInterval, testServer.URL)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := agent.sendMetric(test.args.mType, test.args.name, test.args.value)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
