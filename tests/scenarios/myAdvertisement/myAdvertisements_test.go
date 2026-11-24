package myAdvertisement

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// Positive: Проверяем получение собственных объявлений в профиле
func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	t := s.T()
	title := "My advertisement test " + utils.RandomString(10)

	createBody := advertisement.CreateAdvertisement(t, s.loginData.Token, title, "My advertisement test description", 123.45, 7, filepath.Clean("../testdata/test.jpg"))

	adID := gjson.Get(createBody, "id").String()
	require.NotEmpty(t, adID, "ID созданного объявления не должен быть пустым")

	t.Cleanup(func() {
		advertisement.DeleteAdvertisement(t, s.loginData.Token, adID)
	})

	base.Precondition("Получение списка собственных объявлений")
	advertisementsBody := myAdvertisements.GetMyAdvertisements(t, s.loginData.Token, http.StatusOK, 50)
	items := gjson.Get(advertisementsBody, "items").Array()

	var createdAdvertisement gjson.Result
	found := false

	for _, item := range items {
		if item.Get("id").String() == adID {
			createdAdvertisement = item
			found = true
			break
		}
	}

	require.True(t, found, "Созданное объявление отсутствует в списке собственных объявлений")

	item := createdAdvertisement.String()
	require.Equal(t, s.loginData.User.Id, gjson.Get(item, "user_id").String(), "user_id не совпадает")
	require.NotEmpty(t, gjson.Get(item, "id").String(), "id пустое")
	require.NotEmpty(t, gjson.Get(item, "created_at").String(), "created_at пустое")
	require.NotEmpty(t, gjson.Get(item, "updated_at").String(), "updated_at пустое")
	require.Equal(t, title, gjson.Get(item, "title").String(), "Название объявления не совпадает")
	require.NotEmpty(t, gjson.Get(item, "description").String(), "Описание пустое")
	require.NotEmpty(t, gjson.Get(item, "photos").Array(), "У объявления отсутствуют фотографии")
}

// Negative: Проверяем 401 ошибку при попытке получить собственные объявления без токена
func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений с неправильным авторизационным токеном", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), "incorrect_token", http.StatusUnauthorized, 50)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		require.Equal(s.T(), "unauthorized", gjson.Get(advertisementsBody, "error").String(),
			"Некорректное значение error параметра")
		require.Equal(s.T(), "Invalid or expired token", gjson.Get(advertisementsBody, "message").String(),
			"Некорректное значение message параметра")
	})
}
