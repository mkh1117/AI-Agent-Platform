package httpserver

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {

	engine *gin.Engine

	addr   string
}

func New(addr string, db *gorm.DB, logger *slog.Logger) (*Server,error){
	sqlDB,err := db.DB()

	if err != nil {
		return nil,err
	}

	engine := gin.New()
	
	engine.Use(
		RequestID(),
		RequestLogger(logger),
		Recovery(logger),
	)

	engine.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	engine.GET("/health/ready",func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			2*time.Second,
		)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"checks": gin.H{
					"postgres": "unavailable",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"checks": gin.H{
				"postgres": "ok",
			},
		})
	})

	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":      "not_found",
				"message":   "resource not found",
				"request_id": c.GetString("request_id"),
			},
		})
	})

	return &Server{
		engine: engine,
		addr:   addr,
	}, nil
}

func (s *Server) Handler() http.Handler {
	return s.engine
}

func (s *Server) Addr() string {
	return s.addr
}

func  RequestID() gin.HandlerFunc{
	return func(c *gin.Context) {

		id := make([]byte,16)

		if _, err := rand.Read(id); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "internal_error",
					"message": "internal server error",
				},
			})
			return
		}

		requestID := hex.EncodeToString(id)

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()

	}
}

func RequestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		logger.Info(
			"http request",
			"request_id", c.GetString("request_id"),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
			"client_ip", c.ClientIP(),
		)

	}
}

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error(
					"panic recovered",
					"request_id", c.GetString("request_id"),
					"panic", recovered,
				)

				if !c.Writer.Written() {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": gin.H{
							"code":       "internal_error",
							"message":    "internal server error",
							"request_id": c.GetString("request_id"),
						},
					})
				}

				c.Abort()
			}
		}()

		c.Next()
	}
}