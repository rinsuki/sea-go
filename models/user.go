package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	ID                uint       `json:"id"`
	Name              string     `json:"name"`
	ScreenName        string     `json:"screenName"`
	EncryptedPassword string     `json:"-"`
	PostsCount        int        `json:"postsCount"`
	InviteCodeID      *uint      `json:"-"`
	CanMakeInviteCode bool       `json:"-"`
	AvatarFile        *AlbumFile `json:"avatarFile"`
	AvatarFileID      *uint      `json:"-"`
	MinReadableDate   time.Time  `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (user *User) Pack(db *gorm.DB) {
	if user.AvatarFileID != nil {
		var file AlbumFile
		if err := db.Where("id = ?", user.AvatarFileID).First(&file).Error; err != nil {
			panic(err)
		}
		file.Pack(db)
		user.AvatarFile = &file
	}
}

type Users []User

func (users Users) Pack(db *gorm.DB) {
	var avatarIds []uint
	for _, u := range users {
		if u.AvatarFileID != nil {
			avatarIds = append(avatarIds, *u.AvatarFileID)
		}
	}

	var files AlbumFiles

	if err := db.Where("id IN (?)", avatarIds).Find(&files).Error; err != nil {
		panic(err)
	}

	files.Pack(db)

	for i, u := range users {
		if u.AvatarFileID == nil {
			continue
		}
		var file *AlbumFile
		for _, f := range files {
			if f.ID == *u.AvatarFileID {
				file = &f
				break
			}
		}
		users[i].AvatarFile = file
	}
}
