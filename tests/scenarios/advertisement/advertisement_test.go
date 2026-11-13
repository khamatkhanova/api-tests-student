package advertisement

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	advertisementManager "api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

const (
	testAdvertisementDescription = "Test advertisement description"
	testAdvertisementPrice       = 123.45
	testAdvertisementQuantity    = 7
	testAdvertisementPhoto       = "../../testdata/test.jpg"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Авторизация пользователя для проверки создания объявления")

	s.loginData = auth.Login(
		s.T(),
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// TestCreateAdvertisementPositive проверяет создание объявления и доступность данных
func (s *TestSuite) TestCreateAdvertisementPositive() {
	t := s.T()

	photoPath := filepath.Clean(testAdvertisementPhoto)
	_, err := os.Stat(photoPath)
	require.NoError(t, err, "Не найдена тестовая фотография: "+photoPath)

	title := "API test " + utils.RandomString(10)

	base.Precondition("Создание объявления со всеми параметрами и фотографией")

	createBody := advertisementManager.CreateAdvertisement(t, s.loginData.Token, title, testAdvertisementDescription, testAdvertisementPrice, testAdvertisementQuantity, photoPath)

	id := gjson.Get(createBody, "id").String()
	require.NotEmpty(t, id, "В ответе создания объявления отсутствует ID")

	t.Cleanup(func() {
		advertisementManager.DeleteAdvertisement(t, s.loginData.Token, id)
	})

	base.Precondition("Получение созданного объявления по ID")

	getBody := advertisementManager.GetAdvertisement(t, id)

	require.Equal(t, title, gjson.Get(getBody, "title").String())
	require.Equal(t, testAdvertisementDescription, gjson.Get(getBody, "description").String())
	require.Equal(t, testAdvertisementPrice, gjson.Get(getBody, "price").Float())
	require.Equal(t, testAdvertisementQuantity, int(gjson.Get(getBody, "quantity").Int()))
	require.Equal(t, s.loginData.User.Id, gjson.Get(getBody, "user_id").String())

	base.Precondition("Поиск объявления по полному названию")

	searchBody := advertisementManager.GetAdvertisements(t, title)
	items := gjson.Get(searchBody, "items").Array()

	found := false
	for _, item := range items {
		if gjson.Get(item.String(), "id").String() == id {
			found = true
			break
		}
	}

	require.True(t, found, "Созданное объявление не найдено по полному названию")

	base.Precondition("Проверка загруженной фотографии объявления")

	photosBody := advertisementManager.GetAdvertisementPhotos(t, s.loginData.Token, id)
	photos := gjson.Parse(photosBody).Array()

	require.NotEmpty(t, photos, "У объявления отсутствуют загруженные фотографии")
	require.NotEmpty(t, gjson.Get(photos[0].String(), "url").String(), "У фотографии отсутствует URL")
}

// проверяет создание объявления без авторизации
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	t := s.T()
	photoPath := filepath.Clean(testAdvertisementPhoto)

	testCases := []struct {
		name  string
		token string
	}{
		{name: "Без токена", token: ""},
		{name: "С неправильным токеном", token: "incorrect_token"},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			base.Precondition("Создание объявления: " + testCase.name)
			body := advertisementManager.CreateAdvertisementWithStatus(t, testCase.token, "Unauthorized test", testAdvertisementDescription, testAdvertisementPrice, testAdvertisementQuantity, photoPath, 401)
			require.Equal(t, "unauthorized", gjson.Get(body, "error").String())
		})
	}
}

// TestCreateAdvertisementNegative проверяет создание объявления с некорректными параметрами
func (s *TestSuite) TestCreateAdvertisementNegative() {
	testCases := []struct {
		name        string
		title       string
		description string
		price       float64
	}{
		{name: "Пустое название", title: "", description: testAdvertisementDescription, price: testAdvertisementPrice},
		{name: "Отрицательная цена", title: "Negative price test", description: testAdvertisementDescription, price: -100},
		{name: "Слишком длинное название", title: strings.Repeat("A", 256), description: testAdvertisementDescription, price: testAdvertisementPrice},
		{name: "Пустое описание", title: "Empty description test", description: "", price: testAdvertisementPrice},
	}

	for _, testCase := range testCases {
		s.Run(testCase.name, func() {
			t := s.T()
			base.Precondition("Проверка: " + testCase.name)
			body := advertisementManager.CreateAdvertisementWithStatus(t, s.loginData.Token, testCase.title, testCase.description, testCase.price, testAdvertisementQuantity, filepath.Clean(testAdvertisementPhoto), 400)
			require.NotEmpty(t, body, "Ответ API не должен быть пустым")
		})
	}
}
