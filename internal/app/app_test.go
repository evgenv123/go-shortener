package app

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMyHandler(t *testing.T) {
	type input struct {
		uri    string
		method string
		body   string
	}
	type output struct {
		code     int
		response string
	}

	tests := []struct {
		name string
		inp  input
		outp output
	}{
		// определяем все тесты
		{
			name: "Test #1 (GET)",
			inp: input{
				uri:    "/",
				method: http.MethodGet,
				body:   "",
			},
			outp: output{
				code:     http.StatusBadRequest,
				response: "Wrong request!\n",
			},
		},
		{
			name: "Test #2 (POST)",
			inp: input{
				uri:    "/",
				method: http.MethodPost,
				body:   "https://mail.ru",
			},
			outp: output{
				code: http.StatusCreated,
				// response: "any",
			},
		},
		{
			name: "Test #3 (POST wrong body)",
			inp: input{
				uri:    "/asd",
				method: http.MethodPost,
				body:   "https://yandex.ru",
			},
			outp: output{
				code: http.StatusBadRequest,
				// response: "any",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.inp.method, tt.inp.uri, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(MyHandler)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, res.StatusCode, tt.outp.code, "Wrong status code")

			// тело запроса
			defer res.Body.Close()
			// just adding line
			resBody, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err, "Fail reading body")
			if tt.outp.response != "" {
				assert.Equal(t, string(resBody), tt.outp.response, "Wrong body received")
			}
		})
	}
}
