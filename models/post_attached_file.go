package models

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type PostAttachedFile struct {
	PostID      uint
	AlbumFileID uint
	Order       uint

	AlbumFile AlbumFile
}

func (f PostAttachedFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.AlbumFile)
}

func (PostAttachedFile) TableName() string {
	return "posts_attached_files"
}

type PostAttachedFiles []PostAttachedFile

func (files PostAttachedFiles) Pack(db *gorm.DB) {
	fileIds := make([]uint, len(files))

	for i, f := range files {
		fileIds[i] = f.AlbumFileID
	}

	var albumFiles AlbumFiles
	if err := db.Where("id IN (?)", fileIds).Find(&albumFiles).Error; err != nil {
		panic(err)
	}

	albumFiles.Pack(db)

	for i, f := range files {
		for _, af := range albumFiles {
			if af.ID == f.AlbumFileID {
				files[i].AlbumFile = af
			}
		}
	}
}
