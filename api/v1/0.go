package v1

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rinsuki/sea-go/db"
	"github.com/rinsuki/sea-go/models"
)

func getToken(c *gin.Context) models.AccessToken {
	return c.MustGet("token").(models.AccessToken)
}

func getUser(c *gin.Context) models.User {
	return c.MustGet("user").(models.User)
}

func apiError(c *gin.Context, status int, errors ...string) {
	errorObjects := make([]gin.H, len(errors))

	for i, err := range errors {
		errorObjects[i] = gin.H{
			"message": err,
		}
	}

	c.JSON(status, gin.H{
		"errors": errorObjects,
	})
	c.Abort()
}

func RegisterToRouter(router *gin.RouterGroup) {
	db := db.GetConnection()

	router.Use(func(c *gin.Context) { // CORS support
		if c.GetHeader("origin") == "" { // Originがないならどうでもいい
			return
		}
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	})

	router.OPTIONS("*path", func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, PATCH")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Add("Access-Control-Max-Age", "86400" /* 1 day as a seconds*/)
	})

	router.Use(func(c *gin.Context) { // Error handling
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			apiError(c, 503, "Internal Server Error")
			panic(err)
		}()

		c.Next()
	})

	router.Use(func(c *gin.Context) { // Authorization

		token := c.GetHeader("Authorization")
		if token == "" {
			apiError(c, 400, "Please authorize")
			return
		}

		tokenFields := strings.Fields(token)

		if len(tokenFields) != 2 {
			apiError(c, 400, "Invalid authorize format")
			return
		}

		if tokenFields[0] != "Bearer" {
			apiError(c, 400, "Authorize type is invalid")
			return
		}

		token = tokenFields[1]

		var t models.AccessToken
		if err := db.Where("token = ?", token).First(&t).Error; err != nil {
			if err.Error() == "record not found" {
				apiError(c, 400, "Authorize failed")
				return
			}
			panic(err)
		}

		if t.RevokedAt.Valid {
			apiError(c, 403, "This token is already revoked")
			return
		}

		var user models.User

		if err := db.Where("id = ?", t.UserID).Find(&user).Error; err != nil {
			panic(err)
		}

		if user.InviteCodeID == nil {
			apiError(c, 400, "Please check web interface")
			return
		}

		c.Set("token", t)
		c.Set("user", user)
	})

	router.GET("/account", getAuthorizedAccount)
	router.PATCH("/account", updateAccountProfile)

	router.GET("/timelines/public", getPublicTimeline)
}
