package handlers

import (
	"net/http"

	rsp "github.com/SyntSugar/ss-infra-go/api/response"
	"github.com/gin-gonic/gin"
)

func NoRoute(c *gin.Context) {
	rsp.ResponseWithErrors(c, http.StatusNotFound, 0, nil)
}

func NoMethod(c *gin.Context) {
	rsp.ResponseWithErrors(c, http.StatusMethodNotAllowed, 0, nil)
}
