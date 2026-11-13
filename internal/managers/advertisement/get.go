package advertisement

import (
	"io"
	"net/http"
	"testing"

	httpAdvertisement "api-tests-template/internal/client/http/advertisement"

	"github.com/stretchr/testify/require"
)

// проверяет статус ответа и возвращает его тело
func readResponseBody(t *testing.T, response *http.Response, expectedStatus int) string {
	t.Helper()

	require.Equal(t, expectedStatus, response.StatusCode, "Некорректный HTTP-статус")

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать тело ответа")
	return string(body)
}

// GetAdvertisement получает объявление по ID
func GetAdvertisement(t *testing.T, id string) string {
	response := httpAdvertisement.HttpGetAdvertisement(t, id)
	return readResponseBody(t, response, http.StatusOK)
}

// GetAdvertisements ищет объявления по названию
func GetAdvertisements(t *testing.T, title string) string {
	response := httpAdvertisement.HttpGetAdvertisements(t, title)
	return readResponseBody(t, response, http.StatusOK)
}

// GetAdvertisementPhotos получает фотографии объявления
func GetAdvertisementPhotos(t *testing.T, token, id string) string {
	response := httpAdvertisement.HttpGetAdvertisementPhotos(t, token, id)
	return readResponseBody(t, response, http.StatusOK)
}

// DeleteAdvertisement удаляет объявление по ID
func DeleteAdvertisement(t *testing.T, token, id string) {
	response := httpAdvertisement.HttpDeleteAdvertisement(t, token, id)
	readResponseBody(t, response, http.StatusNoContent)
}
