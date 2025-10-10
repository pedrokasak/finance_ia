package router

import (
	"log"
	"os"
	"strings"

	"finance-ia/internal/config/middleware"

	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
    RegisterRoutes(public, protected gin.IRouter)
}

func NewRouter(registrars ...RouteRegistrar) *gin.Engine {
    r := gin.New()

    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "debug"
    }
    switch env {
    case "debug", "test":
        gin.SetMode(gin.DebugMode)
    case "release":
        gin.SetMode(gin.ReleaseMode)
    }

    tp := os.Getenv("TRUSTED_PROXIES")
    if tp == "" {
        if err := r.SetTrustedProxies([]string{"127.0.0.1", "::1"}); err != nil {
            log.Printf("falha ao definir trusted proxies padrao: %v", err)
        }
    } else {
        proxies := strings.Split(tp, ",")
        if err := r.SetTrustedProxies(proxies); err != nil {
            log.Printf("falha ao definir trusted proxies a partir de TRUSTED_PROXIES: %v", err)
        }
    }
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.Use(middleware.CORS())

    public := r.Group("/api/v1")
    protected := r.Group("/api/v1", middleware.JWTAuth())

    for _, reg := range registrars {
        reg.RegisterRoutes(public, protected)
    }

    public.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

    return r
}