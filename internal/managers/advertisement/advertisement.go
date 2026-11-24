package advertisement

import (
	"io"
	"net/http"
	"testing"

	httpAdvertisement "api-tests-template/internal/client/http/advertisement"

	"github.com/stretchr/testify/require"
)

// CreateAdvertisement создаёт объявление и возвращает ответ
func CreateAdvertisement(t *testing.T, token string, title string, description string, price float64, quantity int, photoPath string) string {
	response, err := httpAdvertisement.HttpPostAdvertisement(token, title, description, price, quantity, photoPath)
	require.NoError(t, err, "Не удалось отправить запрос на создание объявления")
	defer response.Body.Close()

	require.Equal(t, http.StatusCreated, response.StatusCode, "Ожидался статус 201 при создании объявления")

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ сервера")

	return string(responseBody)
}

// CreateAdvertisementWithStatus создаёт объявление и проверяет ожидаемый статус
func CreateAdvertisementWithStatus(t *testing.T, token string, title string, description string, price float64, quantity int, photoPath string, expectedStatus int) string {
	response, err := httpAdvertisement.HttpPostAdvertisement(token, title, description, price, quantity, photoPath)
	require.NoError(t, err, "Не удалось отправить запрос на создание объявления")
	defer response.Body.Close()

	require.Equal(t, expectedStatus, response.StatusCode,
		"Некорректный HTTP статус при создании объявления")

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ сервера")

	return string(responseBody)
}

// CreateAdvertisementWithParams отправляет объявление с произвольными параметрами и проверяет ожидаемый статус
func CreateAdvertisementWithParams(t *testing.T, token, title, description, price, quantity string, photoPaths []string, expectedStatus int) string {
	t.Helper()

	response, err := httpAdvertisement.HttpPostAdvertisementWithParams(token, title, description, price, quantity, photoPaths)
	require.NoError(t, err, "Не удалось отправить запрос на создание объявления")
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ сервера")
	require.Equal(t, expectedStatus, response.StatusCode, "Некорректный HTTP-статус при создании объявления: "+string(responseBody))

	return string(responseBody)
}
