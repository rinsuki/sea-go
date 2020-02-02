package models

import "github.com/jinzhu/gorm"

type AlbumFile struct {
	ID       uint               `json:"id"`
	UserID   uint               `json:"-"`
	Name     string             `json:"name"`
	Type     string             `json:"type"`
	Varaints []AlbumFileVariant `json:"variants"`
}

func (file *AlbumFile) Pack(db *gorm.DB) {
	if err := db.Where("deleted_at IS NULL").Where("album_file_id = ?", file.ID).Order("score DESC").Find(&file.Varaints).Error; err != nil {
		panic(err)
	}
}

type AlbumFiles []AlbumFile

func (files AlbumFiles) Pack(db *gorm.DB) {
	fileIds := make([]uint, len(files))
	for i, f := range files {
		fileIds[i] = f.ID
	}

	var fileVariants []AlbumFileVariant

	if err := db.Where("deleted_at IS NULL").Where("album_file_id IN (?)", fileIds).Order("score DESC").Find(&fileVariants).Error; err != nil {
		panic(err)
	}
	for i, file := range files {
		for _, v := range fileVariants {
			if v.AlbumFileID == file.ID {
				file.Varaints = append(file.Varaints, v)
			}
		}
		files[i] = file
	}
}
