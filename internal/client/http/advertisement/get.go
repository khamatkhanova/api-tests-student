package advertisement

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetAdvertisement получает объявление по ID
func HttpGetAdvertisement(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisements ищет объявления по названию
func HttpGetAdvertisements(t *testing.T, title string) *http.Response {
	return apiRunner.GetRunner().Create().
		Get(path.AdvertisementsPath).
		Query("search", title).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos получает фотографии объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	photosPath := fmt.Sprintf(path.AdvertisementPhotosPath, id)

	return apiRunner.GetRunner().Auth(token).Create().
		Get(photosPath).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement удаляет объявление по ID
func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
