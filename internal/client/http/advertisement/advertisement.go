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

// HttpPostAdvertisement отправляет запрос на создание объявления с фотографией
func HttpPostAdvertisement(token string, title string, description string, price float64, quantity int, photoPath string) (*http.Response, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fields := map[string]string{
		"title":       title,
		"description": description,
		"price":       strconv.FormatFloat(price, 'f', -1, 64),
		"quantity":    strconv.Itoa(quantity),
	}

	for name, value := range fields {
		if err := writer.WriteField(name, value); err != nil {
			return nil, err
		}
	}

	file, err := os.Open(photoPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("photos", filepath.Base(photoPath))
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, err
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
