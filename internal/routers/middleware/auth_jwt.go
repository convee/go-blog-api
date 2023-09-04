package middleware

import (
	"github.com/convee/go-blog-api/internal/enum"
	"github.com/convee/go-blog-api/internal/models"
	"github.com/convee/go-blog-api/internal/pkg/app"
	"github.com/convee/go-blog-api/internal/pkg/e"
	"github.com/convee/go-blog-api/pkg/jwt"
	"github.com/convee/go-blog-api/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type JWTHeader struct {
	Authorization string `header:"Authorization" validate:"required"`
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		var ag = app.Gin{C: c}

		code = e.SUCCESS
		var headers JWTHeader
		validateErr := app.BindHeader(c, &headers)
		if len(validateErr) > 0 {
			ag.Error(e.INVALID_PARAMS, "", validateErr)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(headers.Authorization, "Bearer ")
		if token == "" {
			code = e.INVALID_PARAMS
			c.Next()
			return
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			code = e.TOKEN_INVALID
			logger.Error("token err", zap.String("token", token), zap.Error(err))
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.TOKEN_INVALID
			logger.Error("token expire", zap.String("token", token))
		}
		if code != e.SUCCESS {
			ag.Response(http.StatusUnauthorized, code, data)
			c.Abort()
			return
		}
		authInfo := &models.AuthInfo{
			Username: claims.Subject,
		}
		c.Set(enum.UserValueAuth, authInfo)
		c.Next()
	}
}
