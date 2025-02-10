package router

import (
	"encoding/json"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	ginMiddleware "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"

	"newspaper-api/adapter/controller/gin/handler"
	"newspaper-api/adapter/controller/gin/middleware"
	"newspaper-api/adapter/controller/gin/presenter"
	"newspaper-api/adapter/gateway"
	"newspaper-api/pkg"
	"newspaper-api/pkg/logger"
	"newspaper-api/usecase"
)

// Swaggerの設定をする
func setupSwagger(router *gin.Engine) (*openapi3.T, error) {
	swagger, err := presenter.GetSwagger()
	if err != nil {
		return nil, err
	}

	env := pkg.GetEnvDefault("APP_ENV", "development")
	if env == "development" {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	return swagger, nil
}

func NewGinRouter(db *gorm.DB, corsAllowOrigins []string) (*gin.Engine, error) {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware(corsAllowOrigins))
	swagger, err := setupSwagger(router)
	if err != nil {
		logger.Warn(err.Error())
		return nil, err
	}

	router.Use(middleware.GinZap())
	router.Use(middleware.RecoveryWithZap())

	// ViewのHTMLの設定
	router.LoadHTMLGlob("./adapter/presenter/html/*")
	router.GET("/", handler.Index)

	// Healthチェック用のAPI
	router.GET("/health", handler.Health)

	apiGroup := router.Group("/api")
	{
		apiGroup.Use(middleware.TimeoutMiddleware(2 * time.Second))
		v1 := apiGroup.Group("/v1")
		{
			v1.Use(ginMiddleware.OapiRequestValidator(swagger))
			// Newspaper APIを追加
			newspaperRepository := gateway.NewNewspaperRepository(db)
			newspaperUseCase := usecase.NewNewspaperUseCase(newspaperRepository)
			newspaperHandler := handler.NewNewspaperHandler(newspaperUseCase)
			presenter.RegisterHandlers(v1, newspaperHandler)
		}
	}
	return router, err
}
