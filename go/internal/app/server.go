package app

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"taoshop/internal/cache"
	"taoshop/internal/config"
	"taoshop/internal/database"
	"taoshop/internal/http/handler"
	"taoshop/internal/http/middleware"
	"taoshop/internal/repository"
	"taoshop/internal/service"
)

type Server struct {
	engine *gin.Engine
	cfg    config.Config
	db     *sql.DB
	redis  *redis.Client
}

func NewServer() (*Server, error) {
	cfg := config.Load()

	// 1. 先连基础设施，保证应用启动时就能发现配置或依赖问题。
	db, err := database.NewMySQL(cfg.MySQLDSN)
	if err != nil {
		return nil, err
	}

	// 2. 启动时自动建表和初始化商品，降低第一次运行门槛。
	if err := database.EnsureSchema(db); err != nil {
		return nil, err
	}

	redisClient, err := cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Printf("redis unavailable, fallback without cache/rate-limit: %v", err)
		redisClient = nil
	}

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	catalogService := service.NewCatalogService(productRepo, redisClient)
	orderService := service.NewOrderService(db, productRepo, orderRepo, catalogService)

	authHandler := handler.NewAuthHandler(authService)
	catalogHandler := handler.NewCatalogHandler(catalogService)
	orderHandler := handler.NewOrderHandler(orderService)

	// 3. Gin 负责接 HTTP 请求，把请求交给 handler，再进入 service/repository。
	engine := gin.Default()
	engine.Use(cors(cfg.FrontendOrigin))
	engine.Use(middleware.RateLimit(redisClient, 120, time.Minute))
	engine.Static("/assets", "./web/assets")
	engine.StaticFile("/", "./web/index.html")

	api := engine.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/products", catalogHandler.ListProducts)

		protected := api.Group("")
		protected.Use(middleware.Auth(authService))
		protected.POST("/orders", orderHandler.Create)
		protected.GET("/orders/me", orderHandler.Mine)
	}

	return &Server{
		engine: engine,
		cfg:    cfg,
		db:     db,
		redis:  redisClient,
	}, nil
}

func (s *Server) Run() error {
	return s.engine.Run(s.cfg.Addr())
}
