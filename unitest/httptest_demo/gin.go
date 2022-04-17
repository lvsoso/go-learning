package httptest_demo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Param struct {
	Name string `json:"name"`
}

func helloHandler(c *gin.Context) {
	var p Param
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "we need a name",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("hello %s", p.Name),
	})
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/hello", helloHandler)
	return router
}
