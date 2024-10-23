package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/imirjar/rb-diver/config"
	"github.com/imirjar/rb-diver/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	cfg = &config.Config{
		Secret:  "",
		Michman: "localhost:9090",
	}
	fs = &FakeService{}
)

func TestReportsListHandler(t *testing.T) {

	app := &HTTP{
		Service: fs,
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	type sent struct {
		method      string
		body        string
		path        string
		contentType string
	}
	tests := []struct {
		name string
		sent sent
		want want
	}{
		{
			name: "positive test #1",
			sent: sent{
				method:      "POST",
				body:        ``,
				path:        "/reports/list",
				contentType: "application/json",
			},

			want: want{
				code: 200,
				response: `[
					{
						"id": "2",
						"name": "metrics id",
						"query": "SELECT id FROM metrics;"
					}
				]`,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(test.sent.body)
			if err != nil {
				log.Print(err)
			}
			request := httptest.NewRequest(test.sent.method, test.sent.path, &buf)
			// создаём новый Recorder

			w := httptest.NewRecorder()
			app.ReportsListHandler(w, request)

			res := w.Result()
			// проверяем код ответа
			fmt.Println(test.want.code, res.StatusCode, "####")
			assert.Equal(t, test.want.code, res.StatusCode)
			fmt.Println(test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestExecuteHandler(t *testing.T) {

	app := &HTTP{
		Service: fs,
	}
	type want struct {
		code        int
		contentType string
	}
	type sent struct {
		method      string
		body        models.Report
		path        string
		contentType string
	}
	tests := []struct {
		name string
		sent sent
		want want
	}{
		{
			name: "positive test #1",
			sent: sent{
				method:      "POST",
				body:        models.Report{Id: "2"},
				path:        "/execute/",
				contentType: "application/json",
			},

			want: want{
				code:        200,
				contentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(test.sent.body)
			if err != nil {
				log.Print(err)
			}
			request := httptest.NewRequest(test.sent.method, test.sent.path, &buf)

			w := httptest.NewRecorder()
			app.ExecuteHandler(w, request)

			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			require.NoError(t, err)
		})
	}
}

type FakeService struct{}

func (fs *FakeService) Execute(ctx context.Context, id string) (string, error) {
	return "", nil
}
func (fs *FakeService) ReportsList(ctx context.Context) (string, error) {
	return `[
		{
			"id": "2",
			"name": "metrics id",
			"query": "SELECT id FROM metrics;"
		}
	]`, nil
}
