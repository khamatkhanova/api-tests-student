package favorites

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	apiPath "api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpPostFavorite добавляет объявление в избранное
func HttpPostFavorite(token, advertisementID string) (*http.Response, error) {
	body, err := json.Marshal(map[string]string{
		"advertisement_id": advertisementID,
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, apiPath.FavoritesPath, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	return apiRunner.GetRunner().Auth(token).Do(request)
}

// HttpGetFavorites получает список избранного
func HttpGetFavorites(token string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, apiPath.FavoritesPath, nil)
	if err != nil {
		return nil, err
	}

	return apiRunner.GetRunner().Auth(token).Do(request)
}

// HttpDeleteFavorite удаляет объявление из избранного
func HttpDeleteFavorite(token, advertisementID string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(apiPath.FavoritePath, advertisementID), nil)
	if err != nil {
		return nil, err
	}

	return apiRunner.GetRunner().Auth(token).Do(request)
}
