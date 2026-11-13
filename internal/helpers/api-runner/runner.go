package api_runner

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/steinfletcher/apitest"
)

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"

type ApiTest struct {
	token     string
	userAgent string
	host      string
}

func New() *ApiTest {
	return &ApiTest{
		userAgent: defaultUserAgent,
		host:      os.Getenv("API_URL"),
	}
}

func GetRunner() *ApiTest {
	return New()
}

// Auth Установка Bearer token для прохождения авторизации на приватные ручки
func (at *ApiTest) Auth(token string) *ApiTest {
	at.token = token
	return at
}

func (at *ApiTest) Create() *apitest.APITest {
	apitestNew := apitest.New().EnableNetworking()
	if os.Getenv("DEBUG") == "true" {
		apitestNew = apitestNew.Debug()
	}
	return apitestNew.Intercept(func(request *http.Request) {
		request.Header.Set("User-Agent", at.userAgent)
		if err := mergeServiceUrls(at.host, request.URL); err != nil {
			panic(fmt.Sprintf("Не удалось применить хост из API_URL: %s", err.Error()))
		}
		if len(at.token) != 0 {
			request.Header.Add("Authorization", "Bearer "+at.token)
		}
	})
}

func mergeServiceUrls(host string, r *url.URL) error {
	if host == "" {
		return fmt.Errorf("API_URL не задан, проверьте переменные окружения")
	}
	urlParsed, err := url.Parse(host)
	if err != nil {
		return fmt.Errorf("host cannot be parsed: %s", err.Error())
	}
	if urlParsed.Path != "" {
		r.Path = path.Join(urlParsed.Path, r.Path)
	}
	r.Scheme = urlParsed.Scheme
	r.Host = urlParsed.Host
	return nil
}

// Do выполняет HTTP-запрос с настройками API runner
func (at *ApiTest) Do(request *http.Request) (*http.Response, error) {
	if request == nil || request.URL == nil {
		return nil, fmt.Errorf("HTTP запрос или URL не задан")
	}

	if err := mergeServiceUrls(at.host, request.URL); err != nil {
		return nil, fmt.Errorf("не удалось применить хост из API_URL: %w", err)
	}

	if request.URL.Scheme != "http" && request.URL.Scheme != "https" {
		return nil, fmt.Errorf("недопустимая схема URL: %s", request.URL.Scheme)
	}

	apiURL, err := url.Parse(at.host)
	if err != nil {
		return nil, fmt.Errorf("не удалось разобрать API_URL: %w", err)
	}

	if request.URL.Host == "" || request.URL.Host != apiURL.Host {
		return nil, fmt.Errorf("адрес запроса не совпадает с API_URL")
	}

	request.Header.Set("User-Agent", at.userAgent)
	if at.token != "" {
		request.Header.Set("Authorization", "Bearer "+at.token)
	}

	return http.DefaultTransport.RoundTrip(request)
}
