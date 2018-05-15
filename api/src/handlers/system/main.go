package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"fmt"
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("pong %s", time.Now().Format(time.RFC850)))
}