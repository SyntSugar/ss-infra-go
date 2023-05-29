package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponseWithOK(t *testing.T) {
	data := map[string]interface{}{
		"message": "hello, world",
	}
	router := gin.Default()
	router.GET("/ok", func(c *gin.Context) {
		ResponseWithOK(c, data)
	})

	req, _ := http.NewRequest("GET", "/ok", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "hello, world")
}

func TestResponseWithCreated(t *testing.T) {
	data := map[string]interface{}{
		"message": "resource created",
	}
	router := gin.Default()
	router.POST("/create", func(c *gin.Context) {
		ResponseWithCreated(c, data)
	})

	req, _ := http.NewRequest("POST", "/create", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "resource created")
}

func TestResponseWithAccepted(t *testing.T) {
	data := map[string]interface{}{
		"message": "request accepted",
	}
	router := gin.Default()
	router.PUT("/update", func(c *gin.Context) {
		ResponseWithAccepted(c, data)
	})

	req, _ := http.NewRequest("PUT", "/update", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusAccepted, resp.Code)
	assert.Contains(t, resp.Body.String(), "request accepted")
}

func TestResponseWithErrors(t *testing.T) {
	errors := []interface{}{"error occurred"}
	router := gin.Default()
	router.GET("/error", func(c *gin.Context) {
		ResponseWithErrors(c, http.StatusBadRequest, 0, errors)
	})

	req, _ := http.NewRequest("GET", "/error", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "error occurred")
}

func TestUnmarshalResponse(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	response := &Response{
		Meta: Meta{
			Code:    http.StatusOK,
			Type:    "OK",
			Message: "Success",
		},
		Data: data,
	}

	bytes, err := json.Marshal(response)
	assert.Nil(t, err)

	unmarshaledResponse, err := UnmarshalResponse(bytes)
	assert.Nil(t, err)
	assert.Equal(t, response, unmarshaledResponse)
}

func TestGetData(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	response := &Response{
		Meta: Meta{
			Code:    http.StatusOK,
			Type:    "OK",
			Message: "Success",
		},
		Data: data,
	}

	var receivedData map[string]string
	err := response.GetData(&receivedData)
	assert.Nil(t, err)
	assert.Equal(t, data, receivedData)
}

func TestGetDataNonPointer(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	response := &Response{
		Meta: Meta{
			Code:    http.StatusOK,
			Type:    "OK",
			Message: "Success",
		},
		Data: data,
	}

	var receivedData map[string]string
	err := response.GetData(receivedData) // Pass non-pointer
	assert.Equal(t, err, errors.New("unmarshal receiver was non-pointer"))
}
