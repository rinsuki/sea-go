package models

import (
	"encoding/json"
	"fmt"
	"os"
)

type AlbumFileVariant struct {
	ID          uint    `json:"id"`
	AlbumFileID uint    `json:"-"`
	Score       uint    `json:"score"`
	Extension   string  `json:"extension"`
	Type        string  `json:"type"`
	Size        uint    `json:"size"`
	Hash        *string `json:"-"`
}

func (v AlbumFileVariant) GetPublicURL() string {
	if v.Hash == nil { // 古いファイル
		hexID := fmt.Sprintf("%08x", v.AlbumFileID)
		filePrefixPath := "/album_files/" + hexID[0:2] + "/" + hexID[2:4] + "/" + hexID[4:6] + "/" + hexID[6:8]
		return filePrefixPath + "/" + v.Type + "." + v.Extension
	}
	hash := *v.Hash
	return "/files/s5/" + hash[0:2] + "/" + hash[2:4] + "/" + hash[4:]
}

func (v AlbumFileVariant) MarshalJSON() ([]byte, error) {
	type Alias AlbumFileVariant

	return json.Marshal(&struct {
		Alias
		URL  string `json:"url"`
		Mime string `json:"mime"`
	}{
		Alias: (Alias)(v),
		URL:   os.Getenv("S3_PUBLIC_URL") + v.GetPublicURL(),
		Mime: map[string]string{
			"webp": "image/webp",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"mp4":  "video/mp4",
		}[v.Extension],
	})
}
