package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

func buildAPICode(statusCode, subCode int) int {
	code, _ := strconv.Atoi(fmt.Sprintf("%d%0*d", statusCode, 2, subCode))
	return code
}

type Meta struct {
	Code    int    `json:"code"`
	Type    string `json:"type,omitempty"`
	Message string `json:"message,omitempty"`
	Errors  []any  `json:"errors,omitempty"`
}

type Response struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data"`

	bytes []byte
}

// ResponseWithOK would write the success data output with status ok
func ResponseWithOK(ctx *gin.Context, data any) {
	ResponseWithSuccess(ctx, http.StatusOK, data)
}

// ResponseWithCreated would write the success data output with status created
func ResponseWithCreated(ctx *gin.Context, data any) {
	ResponseWithSuccess(ctx, http.StatusCreated, data)
}

// ResponseWithAccepted would write the success data output with status accepted
func ResponseWithAccepted(ctx *gin.Context, data any) {
	ResponseWithSuccess(ctx, http.StatusAccepted, data)
}

// ResponseError would write the error data output with status error
func ResponseError(ctx *gin.Context, statusCode int, err error) {
	ResponseWithErrors(ctx, statusCode, 0, []any{err.Error()})
}

// ResponseWithErrors would write meta with error and empty data
func ResponseWithErrors(ctx *gin.Context, statusCode, subCode int, errors []any) {
	code := buildAPICode(statusCode, subCode)
	desc := HttpCodeDescription(code)

	ctx.JSON(statusCode, &Response{
		Meta: Meta{
			Code:    code,
			Type:    desc.Status,
			Message: desc.Message,
			Errors:  errors,
		},
	})
}

// ResponseWithSuccess would write the success data output and status into http writer
func ResponseWithSuccess(ctx *gin.Context, statusCode int, data any) {
	code := buildAPICode(statusCode, 0)
	desc := HttpCodeDescription(code)
	ctx.JSON(statusCode, &Response{
		Meta: Meta{
			Code:    code,
			Type:    desc.Status,
			Message: desc.Message,
		},
		Data: data,
	})
}

// UnmarshalResponse create response
func UnmarshalResponse(bytes []byte) (*Response, error) {
	rsp := &Response{bytes: bytes}
	if err := json.Unmarshal(bytes, rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}

// GetData was used to unmarshal the data field into obj
func (r *Response) GetData(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("unmarshal receiver was non-pointer")
	}
	return json.Unmarshal(r.bytes, &Response{Data: v})
}

// IsError return whether the response was error or not(http code >= 400)
func (r *Response) IsError() bool {
	return r.Meta.IsError()
}

// Code return code in response's meta
func (r *Response) Code() int {
	return r.Meta.Code
}

// StatusCode return http status code in response's meta
func (r *Response) StatusCode() int {
	return r.Meta.StatusCode()
}

// SubCode return sub code in response's meta
func (r *Response) SubCode() int {
	return r.Meta.SubCode()
}

// Type return type in response's meta
func (r *Response) Type() string {
	return r.Meta.Type
}

// Message return message in response's meta
func (r *Response) Message() string {
	return r.Meta.Message
}

// Errors return errors in response's meta
func (r *Response) Errors() []any {
	return r.Meta.Errors
}

// IsError return whether the response meta was error or not(http code >= 400)
func (m Meta) IsError() bool {
	return m.StatusCode() >= http.StatusBadRequest
}

// StatusCode return http status code in  meta
func (m Meta) StatusCode() int {
	return m.Code / 100
}

// SubCode return sub code in meta
func (m Meta) SubCode() int {
	return m.Code % 100
}
