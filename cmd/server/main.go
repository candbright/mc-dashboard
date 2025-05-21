package main

import (
	"fmt"

	_ "github.com/candbright/mc-dashboard/docs" // 导入swagger文档
	"github.com/candbright/mc-dashboard/internal/app/handler"
	"github.com/candbright/mc-dashboard/internal/app/infra/minecraft"
	repository "github.com/candbright/mc-dashboard/internal/app/infra/respository"
	"github.com/candbright/mc-dashboard/internal/app/usecase"
	"github.com/candbright/mc-dashboard/internal/domain"
	"github.com/candbright/mc-dashboard/internal/pkg/config"
	"github.com/candbright/mc-dashboard/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title MC Dashboard API
// @version 1.0
// @description MC Dashboard 后端服务API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:11223
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请输入Bearer token，格式为：Bearer <token>
func main() {
	// 初始化日志
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatal("Failed to initialize config:", err)
	}

	// 设置gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(config.GlobalConfig.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}
	sqlDB.SetMaxIdleConns(config.GlobalConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.GlobalConfig.Database.MaxOpenConns)

	// 自动迁移数据库表
	if err := db.AutoMigrate(&domain.User{}, &domain.Server{}, &domain.Save{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	//配置mc
	minecraftManager := minecraft.NewManager(minecraft.ManagerConfig{
		RootDir: config.GlobalConfig.MC.Path,
	})

	// 初始化系统模块
	systemHandler := handler.NewSystemHandler()

	// 初始化用户模块
	userRepo := repository.NewUserRepository(db)
	userService := usecase.NewUserService(userRepo, log)
	userHandler := handler.NewUserHandler(userService)

	// 初始化服务器模块
	serverRepo := repository.NewServerRepository(db)
	serverService := usecase.NewServerService(log, minecraftManager, serverRepo)
	serverHandler := handler.NewServerHandler(serverService)

	// 初始化存档模块
	saveRepo := repository.NewSaveRepository(db)
	saveService := usecase.NewSaveService(log, minecraftManager, saveRepo)
	saveHandler := handler.NewSaveHandler(saveService)

	// 初始化路由
	r := gin.Default()

	r.Use(gin.Recovery())

	// 添加错误处理中间件
	r.Use(middleware.ErrorHandler(log))

	// 添加超时中间件
	r.Use(middleware.Timeout(config.GlobalConfig.Server.Timeout))

	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 公开接口
	v1Group := r.Group("/api/v1")
	{
		v1Group.POST("/register", userHandler.Register)
		v1Group.POST("/login", userHandler.Login)
	}

	// 需要认证的接口
	authGroup := v1Group.Group("")
	authGroup.Use(middleware.JWTAuth())
	{
		// 注册系统管理路由
		systemHandler.RegisterRoutes(authGroup)

		// 注册用户管理路由
		userHandler.RegisterRoutes(authGroup)

		// 注册服务器管理路由
		serverHandler.RegisterRoutes(authGroup)

		// 注册存档管理路由
		saveHandler.RegisterRoutes(authGroup)
	}

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
