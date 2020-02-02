package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	ID   uint   `json:"id"`
	Text string `json:"text"`

	User  User              `json:"user"`
	Files PostAttachedFiles `json:"files"`

	UserID        uint  `json:"-"`
	ApplicationID uint  `json:"-"`
	InReplyToID   *uint `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Posts []Post

func (posts Posts) Pack(db *gorm.DB) {
	// pack user
	ids := make([]uint, len(posts))
	userIdsMap := make(map[uint]bool)
	userIds := make([]uint, 0, len(posts))
	for i, p := range posts {
		ids[i] = p.ID
		if userIdsMap[p.UserID] == false {
			userIdsMap[p.UserID] = true
			userIds = append(userIds, p.UserID)
		}
		posts[i].Files = PostAttachedFiles{}
	}

	var users Users
	if err := db.Where("id IN (?)", userIds).Find(&users).Error; err != nil {
		panic(err)
	}
	users.Pack(db)

	// pack attached files

	var attachedFiles PostAttachedFiles
	if err := db.Where("post_id IN (?)", ids).Order("posts_attached_files.order ASC").Find(&attachedFiles).Error; err != nil {
		panic(err)
	}
	attachedFiles.Pack(db)

	for i, p := range posts {
		for _, u := range users {
			if u.ID == p.UserID {
				posts[i].User = u
				break
			}
		}
		for _, f := range attachedFiles {
			if f.PostID == p.ID {
				posts[i].Files = append(posts[i].Files, f)
			}
		}
	}
}
