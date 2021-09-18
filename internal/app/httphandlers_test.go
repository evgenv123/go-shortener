package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMyHandlers(t *testing.T) {
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
			name: "Test GET with no url",
			inp: input{
				uri:    "/",
				method: http.MethodGet,
				body:   "",
			},
			outp: output{
				code: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Test GET with wrong short url",
			inp: input{
				uri:    "/xxx123",
				method: http.MethodGet,
				body:   "",
			},
			outp: output{
				code:     http.StatusBadRequest,
				response: "Wrong requested ID!\n",
			},
		},
		{
			name: "Test POST",
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
			name: "Test POST wrong URL",
			inp: input{
				uri:    "/",
				method: http.MethodPost,
				body:   "https//yandex.ru",
			},
			outp: output{
				code:     http.StatusBadRequest,
				response: "Wrong URL format!\n",
			},
		},
		{
			name: "Test POST new endpoint",
			inp: input{
				uri:    "/api/shorten",
				method: http.MethodPost,
				body:   "{\"url\": \"https://mail.ru\"}",
			},
			outp: output{
				code: http.StatusCreated,
			},
		},
	}
	appConf.BaseURL = "http://localhost:8080"
	appConf.FileStorage = "urlStorage_test.gob"
	appConf.ServerAddr = "localhost:8080"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.inp.method, tt.inp.uri, strings.NewReader(tt.inp.body))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			r := chi.NewRouter()
			// маршрутизация запросов обработчику
			r.Get("/{id}", MyHandlerGetID)
			r.Post("/api/shorten", MyHandlerShorten)
			r.Post("/", MyHandlerPost)
			// запускаем сервер
			r.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, tt.outp.code, res.StatusCode, "Wrong status code")

			// тело запроса
			defer res.Body.Close()

			resBody, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err, "Fail reading body")
			if tt.outp.response != "" {
				assert.Equal(t, tt.outp.response, string(resBody), "Wrong body received")
			}
		})
	}
	assert.NoError(t, os.Remove(appConf.FileStorage), "Cannot remove temp file storage!")
}

// TODO: Include HappyPath to regular tests
func TestHappyPath(t *testing.T) {
	appConf.BaseURL = "http://localhost:8080"
	appConf.FileStorage = "urlStorage_test.gob"
	appConf.ServerAddr = "localhost:8080"

	r := chi.NewRouter()
	// маршрутизация запросов обработчику
	r.Get("/{id}", MyHandlerGetID)
	r.Post("/api/shorten", MyHandlerShorten)
	r.Post("/", MyHandlerPost)
	urlToShorten := "https://mail.ru"

	// создаём новый Recorder
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/", strings.NewReader(urlToShorten))
	// запускаем сервер
	r.ServeHTTP(w, request)
	res := w.Result()
	// проверяем код ответа
	assert.Equal(t, http.StatusCreated, res.StatusCode, "Wrong status code")
	// читаем тело запроса
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err, "Fail reading body")
	// fmt.Println(string(resBody))

	parsedURL := strings.Split(string(resBody), "/")
	assert.GreaterOrEqual(t, len(parsedURL), 4, "Cannot parse body: "+string(resBody))
	// Проверяем обратное преобразование (из сокращенной ссылки)
	reqID := parsedURL[3]
	// создаём новый Recorder
	w2 := httptest.NewRecorder()
	request = httptest.NewRequest("GET", "/"+reqID, nil)
	r.ServeHTTP(w2, request)
	res2 := w2.Result()
	defer res2.Body.Close()
	assert.Equal(t, http.StatusTemporaryRedirect, res2.StatusCode, "Wrong status code")
	unshortenedURL, err := res2.Location()
	assert.NoError(t, err, "Fail reading Location")
	assert.Equal(t, urlToShorten, unshortenedURL.String(), "Wrong unshortened URL!")
	assert.NoError(t, os.Remove(appConf.FileStorage), "Cannot remove temp file storage!")
}
