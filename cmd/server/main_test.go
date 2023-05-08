package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"metric-alert/internal/handlers"
	"metric-alert/internal/storage"
)

func Test_update(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "positive test",
			args: "/update/counter/metric/1",
			want: want{
				code: 200,
			},
		},
		{
			name: "negative test - empty value",
			args: "/update/gauge/metric/",
			want: want{
				code: 400,
			},
		},
		{
			name: "negative test - empty metric name ",
			args: "/update/gauge/",
			want: want{
				code: 404,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ms := storage.NewMemStore()
			request := httptest.NewRequest(http.MethodPost, test.args, nil)
			w := httptest.NewRecorder()
			h := handlers.NewMetric(ms)
			h.UpdateMetric(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, res.StatusCode, test.want.code)
		})
	}
}
