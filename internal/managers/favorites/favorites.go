package favorites

import (
	"io"
	"net/http"
	"testing"

	httpFavorites "api-tests-template/internal/client/http/favorites"

	"github.com/stretchr/testify/require"
)

// readResponseBody проверяет HTTP-статус и возвращает тело ответа
func readResponseBody(t *testing.T, response *http.Response, expectedStatus int) string {
	t.Helper()
	defer response.Body.Close()

	require.Equal(t, expectedStatus, response.StatusCode, "Некорректный HTTP-статус")

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "Не удалось прочитать ответ")

	return string(body)
}

// AddFavorite добавляет объявление в избранное
func AddFavorite(t *testing.T, token, advertisementID string) string {
	return AddFavoriteWithStatus(t, token, advertisementID, http.StatusCreated)
}

// AddFavoriteWithStatus добавляет объявление и проверяет статус
func AddFavoriteWithStatus(t *testing.T, token, advertisementID string, expectedStatus int) string {
	t.Helper()

	response, err := httpFavorites.HttpPostFavorite(token, advertisementID)
	require.NoError(t, err, "Не удалось добавить объявление в избранное")

	return readResponseBody(t, response, expectedStatus)
}

// GetFavorites получает список избранного
func GetFavorites(t *testing.T, token string) string {
	t.Helper()

	response, err := httpFavorites.HttpGetFavorites(token)
	require.NoError(t, err, "Не удалось получить список избранного")

	return readResponseBody(t, response, http.StatusOK)
}

// DeleteFavorite удаляет объявление из избранного
func DeleteFavorite(t *testing.T, token, advertisementID string) {
	t.Helper()

	response, err := httpFavorites.HttpDeleteFavorite(token, advertisementID)
	require.NoError(t, err, "Не удалось удалить объявление из избранного")

	readResponseBody(t, response, http.StatusNoContent)
}

// DeleteFavoriteWithStatus удаляет объявление и проверяет ожидаемый статус
func DeleteFavoriteWithStatus(t *testing.T, token, advertisementID string, expectedStatus int) string {
	t.Helper()

	response, err := httpFavorites.HttpDeleteFavorite(token, advertisementID)
	require.NoError(t, err, "Не удалось отправить запрос на удаление из избранного")

	return readResponseBody(t, response, expectedStatus)
}

// GetFavoritesWithStatus получает список избранного и проверяет ожидаемый статус
func GetFavoritesWithStatus(t *testing.T, token string, expectedStatus int) string {
	t.Helper()

	response, err := httpFavorites.HttpGetFavorites(token)
	require.NoError(t, err, "Не удалось получить список избранного")

	return readResponseBody(t, response, expectedStatus)
}
