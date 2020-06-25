package repository

import (
	"github.com/uwaifo/video_server_api/domain/entity"
)

// PhotoWorkRepositry  . . .
type PhotoWorkRepositry interface {
	SavePhotoWork(*entity.PhotoWork) (*entity.Photowork, map[string]string)
}
