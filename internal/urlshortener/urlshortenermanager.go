package urlshortener

import (
	"errors"
	"math/rand"
	"net/http"

	dcubeerrs "github.com/Imranr2/DCUBE_API/internal/errors"
	"gorm.io/gorm"
)

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const urlLength = 10

type URLShortenerManager interface {
	GetURL(GetRequest) (*GetResponse, dcubeerrs.Error)
	CreateURL(CreateRequest) (*CreateResponse, dcubeerrs.Error)
	DeleteURL(DeleteRequest) (*DeleteResponse, dcubeerrs.Error)
	Redirect(RedirectRequest) (*RedirectResponse, dcubeerrs.Error)
}

type URLShortenerManagerImpl struct {
	database *gorm.DB
}

func NewURLShortenerManager(database *gorm.DB) URLShortenerManager {
	return &URLShortenerManagerImpl{
		database: database,
	}
}

func (m *URLShortenerManagerImpl) GetURL(req GetRequest) (*GetResponse, dcubeerrs.Error) {
	var shortenedURLs []ShortenedURL

	err := m.database.Model(&ShortenedURL{}).Where("user_id = ?", req.UserID).Find(&shortenedURLs).Error

	if err != nil {
		return nil, dcubeerrs.New(http.StatusInternalServerError, "An error occurred while fetching urls")
	}

	if len(shortenedURLs) == 0 {
		return nil, dcubeerrs.New(http.StatusNotFound, "User does not have any URLs")
	}

	return &GetResponse{ShortenedURLs: shortenedURLs}, nil
}

func generateShortenedURL() string {
	b := make([]byte, urlLength)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func (m *URLShortenerManagerImpl) CreateURL(req CreateRequest) (*CreateResponse, dcubeerrs.Error) {
	var shortened string
	var shortenedURL ShortenedURL

	for {
		shortened = generateShortenedURL()
		err := m.database.First(&shortenedURL, ShortenedURL{Shortened: shortened}).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break
			}
			return nil, dcubeerrs.New(http.StatusInternalServerError, "An error occurred while creating shortened url")
		}
	}

	newShortenedURL := ShortenedURL{
		Original:  req.OriginalURL,
		Shortened: shortened,
		UserID:    req.UserID,
	}

	err := m.database.Create(&newShortenedURL).Error

	if err != nil {
		return nil, dcubeerrs.New(http.StatusInternalServerError, "An error occurred while creating shortened url")
	}

	return &CreateResponse{ShortenedURL: newShortenedURL}, nil
}

func (m *URLShortenerManagerImpl) DeleteURL(req DeleteRequest) (*DeleteResponse, dcubeerrs.Error) {
	var shortenedURL ShortenedURL

	err := m.database.First(&shortenedURL, ShortenedURL{ID: req.ID}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dcubeerrs.New(http.StatusNotFound, "URL does not exist")
		}
		return nil, dcubeerrs.New(http.StatusInternalServerError, "An error occurred while deleting url")
	}

	if shortenedURL.UserID != req.UserID {
		return nil, dcubeerrs.New(http.StatusForbidden, "User is trying to delete other users records")
	}

	err = m.database.Delete(&ShortenedURL{}, req.ID).Error

	if err != nil {
		return nil, dcubeerrs.New(http.StatusInternalServerError, "An error occurred while deleting shortened url")
	}

	return &DeleteResponse{ShortenedURL: shortenedURL}, nil
}

func (m *URLShortenerManagerImpl) Redirect(req RedirectRequest) (*RedirectResponse, dcubeerrs.Error) {
	var shortenedURL ShortenedURL

	err := m.database.First(&shortenedURL, ShortenedURL{Shortened: req.URL}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dcubeerrs.New(http.StatusNotFound, "URL does not exist")
		}
		return nil, dcubeerrs.New(http.StatusInternalServerError, "An error occurred while deleting url")
	}

	return &RedirectResponse{OriginalURL: shortenedURL.Original}, nil
}
