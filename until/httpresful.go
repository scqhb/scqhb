package until

import (
	"fmt"
	gin "github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func HttpServer() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	r.GET("/Getfalsepositive/:uuid", func(c *gin.Context) {
		uuidtmp := c.Param("uuid")
		if Filter.TestString(uuidtmp) {
			c.Status(400)
			fmt.Println("############################")
			return
		}

		c.String(http.StatusOK, "false")

		fmt.Println("uuid:::", uuidtmp)
	})

	fmt.Println("启动---9999")

	http.ListenAndServe(":9999", r)

}
