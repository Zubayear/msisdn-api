package main

import (
	"context"
	"fmt"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"huspass/auth"
	"huspass/di_context"
	"huspass/middleware"
	"huspass/model"
	"io"
	"net/http"
	"os"
	"time"
)

var logger *zap.Logger

func main() {

	f, _ := os.OpenFile(fmt.Sprintf("log-%s", time.Now().Format("2006-02-01")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.New()

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{fmt.Sprintf("log-%s", time.Now().Format("2006-02-01"))}
	logger, _ = cfg.Build()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, false))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	//store := persist.NewRedisStore(redis.NewClient(&redis.Options{
	//	Addr:    "127.0.0.1:6379",
	//	Network: "tcp",
	//}))
	err := godotenv.Load()
	if err != nil {
		return
	}

	r.POST("api/signup", createUser)
	r.POST("api/login", getUser)
	r.POST("api/msisdn", middleware.AuthMiddleware() /*middleware.RequestIdMiddleware(),*/, createMsisdn)
	r.GET("api/msisdn/:id", middleware.AuthMiddleware() /*cache.CacheByRequestURI(store, 2*time.Minute),*/ /*middleware.RequestIdMiddleware(),*/, getMsisdnById)
	r.GET("api/msisdn", middleware.AuthMiddleware() /*cache.CacheByRequestURI(store, 2*time.Minute),*/ /*middleware.RequestIdMiddleware(),*/, getMsisdn)
	r.PUT("api/msisdn/:id", middleware.AuthMiddleware() /*middleware.RequestIdMiddleware(),*/, updateMsisdn)
	r.DELETE("api/msisdn/:id", middleware.AuthMiddleware() /*cache.CacheByRequestURI(store, 2*time.Minute),*/ /*middleware.RequestIdMiddleware(),*/, deleteMsisdn)

	port := os.Getenv("SERVERPORT")
	if port == "" {
		port = "8081"
	}
	err = r.Run(":" + port)
	if err != nil {
		return
	}

}

func getUser(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Provide valid parameters",
			"err":     err.Error(),
		})
		return
	}

	dependencyProvider, err := di_context.DependencyProvider()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed creating resource",
			"err":     err.Error(),
		})
		return
	}
	userDto, err := dependencyProvider.UserRepo.GetUser(context.Background(), &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed creating resource",
			"err":     err.Error(),
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDto.Password), []byte(user.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed getting resource",
			"err":     err.Error(),
		})
		return
	}
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"sub": user.Username,
	//	"exp": time.Now().Add(time.Hour).Unix(),
	//})
	//tokenString, err := token.SignedString([]byte("hellotherehowareyouman"))
	token, err := auth.GenerateJWT(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed generating token",
			"err":     err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully generated token",
		"token":   token,
	})
}

func createUser(ctx *gin.Context) {
	var user model.User
	err := ctx.ShouldBindJSON(&user)
	logger.Info("Received request for createUser", zap.Any("request", user))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Provide valid parameters",
			"err":     err.Error(),
		})
		return
	}
	// Save to db
	dependencyProvider, err := di_context.DependencyProvider()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed creating resource",
			"err":     err.Error(),
		})
		return
	}
	err = dependencyProvider.UserRepo.CreateUser(context.Background(), &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed creating resource",
			"err":     err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Successfully created resource",
	})
}

func deleteMsisdn(ctx *gin.Context) {
	dependencyProvider, err := di_context.DependencyProvider()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed deleting resource",
			"err":     err.Error(),
		})
		return
	}
	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed deleting resource",
			"err":     err.Error(),
		})
		return
	}
	err = dependencyProvider.Repo.DeleteMsisdn(context.Background(), parsedId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Failed deleting resource",
			"err":     err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "Successfully deleted resource",
	})
}

func updateMsisdn(ctx *gin.Context) {
	dependencyProvider, err := di_context.DependencyProvider()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed getting resources",
			"err":     err.Error(),
		})
		return
	}

	id := ctx.Param("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed getting resources",
			"err":     err.Error(),
		})
		return
	}
	var msisdn model.Msisdn
	err = ctx.ShouldBindJSON(&msisdn)
	logger.Info("Received request for updateMsisdn", zap.Any("request", msisdn))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed getting resources",
			"err":     err.Error(),
		})
		return
	}

	msisdnDto, err := dependencyProvider.Repo.UpdateMsisdn(context.Background(), &msisdn, parsedId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Failed getting resources",
			"err":     err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully updated resource",
		"data":    msisdnDto,
	})
}

func getMsisdn(ctx *gin.Context) {
	dependencyProvider, err := di_context.DependencyProvider()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed getting resources",
			"err":     err.Error(),
		})
		return
	}

	msisdnDto, err := dependencyProvider.Repo.GetAllMsisdn(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed getting resources",
			"err":     err.Error(),
		})
		return
	}
	//ctx.Writer.Header().Set("Cache-Control", "public, max-age=120000, immutable")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved resources",
		"data":    msisdnDto,
	})
}

func getMsisdnById(ctx *gin.Context) {
	param := ctx.Param("id")
	id, err := uuid.Parse(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Provide valid parameters",
			"err":     err.Error(),
		})
		return
	}
	dependencyProvider, err := di_context.DependencyProvider()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed getting resource",
			"err":     err.Error(),
		})
		return
	}
	msisdnFromRepo, err := dependencyProvider.Repo.GetMsisdn(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Failed getting resource",
			"err":     err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved resource",
		"data":    msisdnFromRepo,
	})

}

func createMsisdn(ctx *gin.Context) {
	var msisdn model.Msisdn
	err := ctx.ShouldBindJSON(&msisdn)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Provide valid parameters",
			"err":     err.Error(),
		})
		return
	}
	// Save to db
	dependencyProvider, err := di_context.DependencyProvider()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed creating resource",
			"err":     err.Error(),
		})
		return
	}
	msisdnDto, err := dependencyProvider.Repo.CreateMsisdn(context.Background(), &msisdn)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed creating resource",
			"err":     err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Successfully created resource",
		"data":    msisdnDto,
	})
}
