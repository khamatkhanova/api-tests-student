package favorites

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	advertisementManager "api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	favoritesManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"
	base "api-tests-template/tests"
)

const (
	favoriteDescription = "Favorite test description"
	favoritePrice       = 150.50
	favoriteQuantity    = 3
	favoritePhoto       = "../testdata/test.jpg"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
	ownerData authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Авторизация пользователей для проверки избранного")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	s.ownerData = auth.Login(s.T(), os.Getenv("TEST_LOGIN_SECOND"), os.Getenv("TEST_PASSWORD_SECOND"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) createTestAdvertisement(t *testing.T, token string) string {
	title := "Favorite test " + utils.RandomString(10)
	body := advertisementManager.CreateAdvertisement(t, token, title, favoriteDescription, favoritePrice, favoriteQuantity, filepath.Clean(favoritePhoto))
	id := gjson.Get(body, "id").String()
	require.NotEmpty(t, id, "ID созданного объявления не должен быть пустым")
	return id
}

func containsFavorite(body, advertisementID string) bool {
	items := gjson.Get(body, "items").Array()
	for _, item := range items {
		if item.Get("advertisement.id").String() == advertisementID {
			return true
		}
	}
	return false
}

// TestFavoriteCycle проверяет добавление и удаление чужого объявления из избранного
func (s *TestSuite) TestFavoriteCycle() {
	t := s.T()
	adID := s.createTestAdvertisement(t, s.loginData.Token)
	t.Cleanup(func() {
		advertisementManager.DeleteAdvertisement(t, s.loginData.Token, adID)
	})

	base.Precondition("Добавление чужого объявления в избранное")
	favoritesManager.AddFavorite(t, s.ownerData.Token, adID)

	base.Precondition("Проверка наличия объявления в избранном")
	body := favoritesManager.GetFavorites(t, s.ownerData.Token)
	require.True(t, containsFavorite(body, adID), "Объявление отсутствует в избранном")

	base.Precondition("Удаление объявления из избранного")
	favoritesManager.DeleteFavorite(t, s.ownerData.Token, adID)

	base.Precondition("Проверка отсутствия объявления в избранном")
	body = favoritesManager.GetFavorites(t, s.ownerData.Token)
	require.False(t, containsFavorite(body, adID), "Объявление осталось в избранном")
}

// TestFavoriteRemovedWithAdvertisement проверяет удаление из избранного после удаления самого объявления
func (s *TestSuite) TestFavoriteRemovedWithAdvertisement() {
	t := s.T()
	adID := s.createTestAdvertisement(t, s.loginData.Token)
	advertisementDeleted := false

	t.Cleanup(func() {
		if !advertisementDeleted {
			advertisementManager.DeleteAdvertisement(t, s.loginData.Token, adID)
		}
	})

	base.Precondition("Добавление чужого объявления в избранное")
	favoritesManager.AddFavorite(t, s.ownerData.Token, adID)

	base.Precondition("Удаление объявления владельцем")
	advertisementManager.DeleteAdvertisement(t, s.loginData.Token, adID)
	advertisementDeleted = true

	base.Precondition("Проверка удаления объявления из избранного")
	body := favoritesManager.GetFavorites(t, s.ownerData.Token)
	require.False(t, containsFavorite(body, adID), "Удалённое объявление осталось в избранном")
}

// TestAddOwnAdvertisementToFavorites проверяет запрет добавления собственного объявления
func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	t := s.T()
	adID := s.createTestAdvertisement(t, s.ownerData.Token)
	t.Cleanup(func() {
		advertisementManager.DeleteAdvertisement(t, s.ownerData.Token, adID)
	})

	base.Precondition("Попытка добавить собственное объявление в избранное")
	body := favoritesManager.AddFavoriteWithStatus(t, s.ownerData.Token, adID, http.StatusForbidden)
	require.NotEmpty(t, body, "Ответ API не должен быть пустым")
}

// TestAddFavoriteTwice проверяет повторное добавление объявления в избранное
func (s *TestSuite) TestAddFavoriteTwice() {
	t := s.T()
	adID := s.createTestAdvertisement(t, s.loginData.Token)

	t.Cleanup(func() {
		advertisementManager.DeleteAdvertisement(t, s.loginData.Token, adID)
	})
	t.Cleanup(func() {
		favoritesManager.DeleteFavoriteWithStatus(t, s.ownerData.Token, adID, http.StatusNoContent)
	})

	base.Precondition("Первое добавление чужого объявления в избранное")
	favoritesManager.AddFavorite(t, s.ownerData.Token, adID)

	base.Precondition("Повторное добавление объявления в избранное")
	body := favoritesManager.AddFavoriteWithStatus(t, s.ownerData.Token, adID, http.StatusConflict)
	require.NotEmpty(t, body, "Ответ API не должен быть пустым")
}

// TestAddFavoriteWithoutAuthorization проверяет запрет добавления в избранное без авторизации
func (s *TestSuite) TestAddFavoriteWithoutAuthorization() {
	t := s.T()
	adID := s.createTestAdvertisement(t, s.loginData.Token)
	t.Cleanup(func() {
		advertisementManager.DeleteAdvertisement(t, s.loginData.Token, adID)
	})

	base.Precondition("Попытка добавить объявление в избранное без авторизации")
	body := favoritesManager.AddFavoriteWithStatus(t, "", adID, http.StatusUnauthorized)
	require.NotEmpty(t, body, "Ответ API не должен быть пустым")
}

// TestAddNonexistentAdvertisementToFavorites проверяет добавление несуществующего объявления
func (s *TestSuite) TestAddNonexistentAdvertisementToFavorites() {
	t := s.T()
	base.Precondition("Попытка добавить несуществующее объявление в избранное")
	body := favoritesManager.AddFavoriteWithStatus(t, s.ownerData.Token, "00000000-0000-0000-0000-000000000000", http.StatusBadRequest)

	require.Equal(t, "invalid_request", gjson.Get(body, "error").String())
	require.Equal(t, "Key: 'AddFavoriteRequest.AdvertisementID' Error:Field validation for 'AdvertisementID' failed on the 'required' tag", gjson.Get(body, "message").String())
}

// TestDeleteNonexistentFavorite проверяет удаление несуществующего объявления из избранного
func (s *TestSuite) TestDeleteNonexistentFavorite() {
	t := s.T()
	base.Precondition("Попытка удалить несуществующее объявление из избранного")
	body := favoritesManager.DeleteFavoriteWithStatus(t, s.ownerData.Token, "00000000-0000-0000-0000-000000000000", http.StatusNotFound)

	require.Equal(t, "not_found", gjson.Get(body, "error").String())
	require.Equal(t, "Favorite entry not found", gjson.Get(body, "message").String())
}

// TestDeleteFavoriteWithoutAuthorization проверяет запрет удаления из избранного без авторизации
func (s *TestSuite) TestDeleteFavoriteWithoutAuthorization() {
	t := s.T()
	adID := s.createTestAdvertisement(t, s.loginData.Token)
	t.Cleanup(func() {
		advertisementManager.DeleteAdvertisement(t, s.loginData.Token, adID)
	})

	base.Precondition("Добавление объявления в избранное")
	favoritesManager.AddFavorite(t, s.ownerData.Token, adID)

	base.Precondition("Попытка удалить объявление из избранного без авторизации")
	body := favoritesManager.DeleteFavoriteWithStatus(t, "", adID, http.StatusUnauthorized)
	require.NotEmpty(t, body, "Ответ API не должен быть пустым")
}

// TestGetFavoritesWithoutAuthorization проверяет запрет получения избранного без авторизации
func (s *TestSuite) TestGetFavoritesWithoutAuthorization() {
	t := s.T()
	base.Precondition("Получение списка избранного без авторизации")
	body := favoritesManager.GetFavoritesWithStatus(t, "", http.StatusUnauthorized)

	require.Equal(t, "unauthorized", gjson.Get(body, "error").String())
}
