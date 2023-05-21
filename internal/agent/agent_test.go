package agent

import (
	"fmt"
	"github.com/rs/zerolog"
	"metric-alert/internal/entities"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestSendMetric(t *testing.T) {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/update" || r.Method != "POST" {
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
		arg  entities.Metrics
		want error
	}{
		{
			name: "positive test",
			arg:  entities.Metrics{MType: "gauge", ID: "metric", Value: &value},
			want: nil,
		},
	}

	reportInterval := 10
	pollInterval := 2
	testURL, err := url.Parse(testServer.URL)
	if err != nil {
		return
	}

	agent := NewAgent(reportInterval, pollInterval, testURL.Host, log)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := agent.sendMetric(test.arg)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
