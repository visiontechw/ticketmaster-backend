package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/visiontechw/ticketmaster/internal/domain/services"
)

func TokenCatcherMiddleware(tokenManager services.TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token de autorização é obrigatório"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := tokenManager.Validate(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		requestCtx := context.WithValue(ctx.Request.Context(), "user_id", claims.UserID)
		ctx.Request = ctx.Request.WithContext(requestCtx)

		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)

		ctx.Next()
	}
}
