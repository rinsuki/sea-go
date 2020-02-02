package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rinsuki/sea-go/db"
	"github.com/rinsuki/sea-go/models"
)

func getPublicTimeline(c *gin.Context) {
	type Params struct {
		SinceID *uint   `form:"sinceId"`
		MaxID   *uint   `form:"maxId"`
		Count   *uint   `form:"count"`
		Search  *string `form:"search"`
	}

	db := db.GetConnection()
	user := getUser(c)

	var params Params
	if err := c.ShouldBind(&params); err != nil {
		panic(err)
	}

	query := db.Where("created_at > ?", user.MinReadableDate).Order("created_at DESC")

	if params.SinceID != nil {
		query = query.Where("id > ?", *params.SinceID)
	}
	if params.MaxID != nil {
		query = query.Where("id < ?", *params.MaxID)
	}
	if params.Search != nil {
		query = query.Where("text LIKE ?", params.Search)
	}
	if params.Count != nil {
		if *params.Count < 1 {
			panic("count too small")
		}
		if *params.Count > 100 {
			panic("count too big")
		}
		query = query.Limit(*params.Count)
	} else {
		query = query.Limit(20)
	}

	var posts models.Posts

	if err := query.Find(&posts).Error; err != nil {
		panic(err)
	}

	posts.Pack(db)

	c.JSON(200, posts)
}
