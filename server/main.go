package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"goST/Controllers"
	"net/http"
)

func main() {
	Controllers.ConnectStripeDB()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Use(static.Serve("/", static.LocalFile("./dist", true)))

	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "api文档")
		})
		api.GET("/dashAccounts", Controllers.GetDashAccounts)
		api.GET("/dashBalances", Controllers.GetDashBalances)
		api.GET("/dashInvoices", Controllers.GetDashInvoices)
		api.GET("/dashPayouts", Controllers.GetDashPayouts)
		api.GET("/updateInfos", Controllers.GetUpdateInfos)
		api.PATCH("/updateInfos", Controllers.PatchUpdateInfos)
	}
	router.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})
	router.Run(":8080")
}
