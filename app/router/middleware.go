package router

import (
	"pm_go_version/app/constant"
	"pm_go_version/app/pkg"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// CORSMiddleware 处理跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的域名列表
		allowedOrigins := []string{
			"http://localhost:3000", // React开发环境
			//"http://localhost:5173",   // Vite开发环境
			//"https://your-domain.com", // 生产环境域名
		}

		origin := c.Request.Header.Get("Origin")
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24小时
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer pkg.PanicHandler(c)

		// 1. 首先检查Authorization header
		token := c.GetHeader("Authorization")
		if len(token) > 7 {
			claims, err := pkg.Verfiy(token[7:])
			if err == nil {
				log.Info("Token verified from Authorization header")
				c.Set("parse_id", claims.ID)
				c.Set("parse_username", claims.UserName)
				c.Set("parse_email", claims.Email)
				c.Next()
				return
			}
		}

		// 2. 如果没有Authorization header或验证失败，检查cookie
		cookieToken, err := c.Cookie("token")
		//log.Info("cookieToken is: ", cookieToken)
		if err == nil && cookieToken != "" {
			claims, err := pkg.Verfiy(cookieToken)
			if err == nil {
				log.Info("Token verified from cookie")
				c.Set("parse_id", claims.ID)
				c.Set("parse_username", claims.UserName)
				c.Set("parse_email", claims.Email)
				c.Next()
				return
			}
		}

		// 3. 如果两种方式都失败，返回未授权错误
		log.Error("Authentication failed: no valid token found in header or cookie")
		pkg.PanicException(constant.Unauthorized)
		c.Abort()
	}
}
