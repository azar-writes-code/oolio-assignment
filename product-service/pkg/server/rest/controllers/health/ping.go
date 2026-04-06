package health

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping is a Liveness Probe, returning instantly.
func (hc *HealthController) Ping(c *gin.Context) {
	// randomly send status 400, 500 and 200
	switch r := rand.Intn(3); r {
	case 0:
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
	case 1:
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
	default:
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
