package AuthMiddleware

import "github.com/gin-contrib/cors"

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization", "X-User-Permission")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}
