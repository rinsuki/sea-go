package v1

import (
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rinsuki/sea-go/db"
	"github.com/rinsuki/sea-go/models"
)

func getAuthorizedAccount(c *gin.Context) {
	db := db.GetConnection()
	user := getUser(c)
	user.Pack(db)

	c.JSON(200, user)
}

func updateAccountProfile(c *gin.Context) {
	omikujiContents := []string{
		"大吉", "中吉", "吉", "小吉", "末吉", "凶", "大凶", "はずれ",
	}

	type Params struct {
		Name         *string `form:"name"`
		AvatarFileId *uint   `form:"avatarFileId"`
	}

	db := db.GetConnection()
	user := getUser(c)

	var params Params
	if err := c.ShouldBind(&params); err != nil {
		panic(err)
	}

	if params.Name != nil {
		name := *params.Name

		if len(name) < 1 {
			panic("name too short")
		} else if len(name) > 20 {
			panic("name too long")
		}

		if name == "!omikuji" { // !!! OMIKUJI MODE!!!
			name += " → ★" + omikujiContents[rand.Intn(len(omikujiContents))]
		} else if strings.Contains(name, "★") {
			name = strings.ReplaceAll(name, "★", "☆")
		}
		user.Name = name
	}

	if params.AvatarFileId != nil {
		if *params.AvatarFileId == 0 { // delete already submitted avatar
			user.AvatarFileID = nil
		} else {
			var file models.AlbumFile
			if err := db.Where("id = ?", params.AvatarFileId).Find(&file).Error; err != nil {
				panic(err)
			}
			if file.UserID != user.ID {
				panic("you are not owner")
			}
			if file.Type != "image" {
				panic("invalid file type")
			}
			user.AvatarFileID = params.AvatarFileId
		}
	}

	if err := db.Model(&user).Select("name", "avatar_file_id").Updates(user).Error; err != nil {
		panic(err)
	}

	user.Pack(db)

	c.JSON(200, user)
}
