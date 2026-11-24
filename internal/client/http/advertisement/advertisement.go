package advertisement

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpPostAdvertisement отправляет запрос на создание объявления с одной фотографией
func HttpPostAdvertisement(token string, title string, description string, price float64, quantity int, photoPath string) (*http.Response, error) {
	return HttpPostAdvertisementWithParams(token, title, description, strconv.FormatFloat(price, 'f', -1, 64), strconv.Itoa(quantity), []string{photoPath})
}

// HttpPostAdvertisementWithParams отправляет запрос с произвольными строковыми параметрами и несколькими фотографиями
func HttpPostAdvertisementWithParams(token, title, description, price, quantity string, photoPaths []string) (*http.Response, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fields := map[string]string{
		"title":       title,
		"description": description,
		"price":       price,
		"quantity":    quantity,
	}

	for name, value := range fields {
		if err := writer.WriteField(name, value); err != nil {
			return nil, err
		}
	}

	for _, photoPath := range photoPaths {
		file, err := os.Open(photoPath)
		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormFile("photos", filepath.Base(photoPath))
		if err != nil {
			file.Close()
			return nil, err
		}

		_, copyErr := io.Copy(part, file)
		closeErr := file.Close()

		if copyErr != nil {
			return nil, copyErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, path.AdvertisementPath, &body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	return apiRunner.GetRunner().Auth(token).Do(request)
}
