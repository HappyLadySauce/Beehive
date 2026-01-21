package middleware

import (
	"github.com/gin-gonic/gin"

	// 注意：需要在 go.mod 中手动添加以下依赖，然后执行 go mod tidy
	// github.com/swaggo/gin-swagger
	// github.com/swaggo/files
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterSwagger 在给定的路由上注册 Swagger UI。
// 约定：各服务需要自行通过 swag 工具生成 docs 包，并配置 docs.SwaggerInfo。
func Swagger() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/v1/swagger/doc.json"))
}
