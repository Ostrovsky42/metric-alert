package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSendMetric(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/update/gauge/metric/1" || r.Method != "POST" {
			t.Errorf("Unexpected request URL or method: %v %v", r.URL.Path, r.Method)
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "text/plain" {
			t.Errorf("Unexpected content type: %v", contentType)
		}
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

	reportInterval := 10
	pollInterval := 2
	testURL, err := url.Parse(testServer.URL)
	if err != nil {
		return
	}

	agent := NewAgent(reportInterval, pollInterval, testURL.Host)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := agent.sendMetric(test.args.mType, test.args.name, test.args.value)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
