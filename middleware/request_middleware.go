package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"huspass/auth"
	"huspass/model"
	"log"
	"net/http"
	"strings"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqId := uuid.New()
		ctx.Set("reqId", reqId)
		log.Printf("Request %s start:", reqId)
		ctx.Next()
		log.Printf("Request %s end", reqId)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "can't access the resource without proper credentials",
			})
			ctx.Abort()
			return
		}
		split := strings.Split(bearerToken, " ")
		token := split[1]
		role, err := auth.ValidateToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			ctx.Abort()
			return
		}
		val := Authorization(role)
		if !val {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "you're not authorized to access the resource",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func Authorization(validRoles model.Role) bool {
	return validRoles == model.ProductManager
}
func Authorize(obj string, act string, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, existed := c.Get("userID")
		if !existed {
			c.AbortWithStatusJSON(401, gin.H{"msg": "User hasn't logged in yet"})
			return
		}
		err := enforcer.LoadPolicy()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"msg": "Failed to load policy from DB"})
			return
		}
		ok, err := enforcer.Enforce(fmt.Sprint(sub), obj, act)

		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"msg": "Error occurred when authorizing user"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"msg": "You are not authorized"})
			return
		}
		c.Next()
	}
}
