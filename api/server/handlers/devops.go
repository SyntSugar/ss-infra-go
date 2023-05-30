package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type status int

var _status status = http.StatusOK

func (s status) String() string {
	return http.StatusText(int(s))
}

func (s *status) Online() {
	*s = http.StatusOK
}

func (s *status) Offline() {
	*s = http.StatusServiceUnavailable
}

func GetDevopsStatus(c *gin.Context) {
	c.String(int(_status), _status.String())
}

func UpdateDevopsStatus(c *gin.Context) {
	if status, ok := c.GetPostForm("status"); ok {
		code, err := strconv.Atoi(status)
		if err != nil || (code != http.StatusOK && code != http.StatusServiceUnavailable) {
			c.String(http.StatusBadRequest, "parameter 'status' should be 200 or 503")
		}
		var msg string
		if code == http.StatusOK {
			_status.Online()
			msg = "the service was online"
		} else {
			_status.Offline()
			msg = "the service was offline"
		}
		c.String(http.StatusOK, msg)
	}
	c.String(http.StatusBadRequest, "parameter 'status' was required")
}
