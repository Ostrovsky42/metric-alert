package agent

import (
	"fmt"
	"metric-alert/internal/agent/gatherer"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSendMetric(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/update/" || r.Method != "POST" {
			t.Errorf("Unexpected request URL or method: %v %v", r.URL.Path, r.Method)
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Unexpected content type: %v", contentType)
		}
		fmt.Fprint(w, "OK")
	}))
	defer testServer.Close()
	var value float64 = 1
	tests := []struct {
		name string
		arg  gatherer.Metrics
		want error
	}{
		{
			name: "positive test",
			arg:  gatherer.Metrics{MType: "gauge", ID: "metric", Value: &value},
			want: nil,
		},
	}

	reportInterval := 10
	pollInterval := 2
	testURL, err := url.Parse(testServer.URL)
	if err != nil {
		return
	}

	agent := NewAgent(reportInterval, pollInterval, testURL.Host)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = agent.sender.SendMetricJSON(test.arg)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
