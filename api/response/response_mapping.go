package response

import "errors"

type Description struct {
	Status  string
	Message string
}

var code2Description = map[int]Description{
	20000: {
		Status:  "OK",
		Message: "The request was successfully processed by AfterShip.",
	},
	20100: {
		Status:  "Created",
		Message: "The request has been fulfilled and a new resource has been created.",
	},
	20200: {
		Status:  "Accepted",
		Message: "The request has been accepted for processing, but the processing has not been completed.",
	},
	40000: {
		Status: "BadRequest",
		Message: "The request was not understood by the server, " +
			"generally due to bad syntax or because the Content-Type header was not correctly set to application/json.",
	},
	40100: {
		Status:  "Unauthorized",
		Message: "The necessary authentication credentials are not present in the request or are incorrect.",
	},
	40101: {
		Status:  "Unauthorized",
		Message: "The necessary authentication credentials are not present in the request or are incorrect.",
	},
	40102: {
		Status:  "Unauthorized",
		Message: "The necessary authentication credentials are not present in the request or are incorrect.",
	},
	40103: {
		Status:  "Unauthorized",
		Message: "The necessary authentication credentials are not present in the request or are incorrect.",
	},
	40200: {
		Status:  "PaymentRequired",
		Message: "Payment Required",
	},
	40300: {
		Status:  "Forbidden",
		Message: "The server is refusing to respond to the request. This is generally because you have not requested the appropriate scope for this action.",
	},
	40400: {
		Status:  "NotFound",
		Message: "The requested resource was not found but could be available again in the future.",
	},
	40500: {
		Status:  "MethodNotAllowed",
		Message: "The method received in the request-line is known by the server but not supported by the target resource.",
	},
	40900: {
		Status:  "Conflict",
		Message: "The request conflicts with another request (perhaps due to using the same idempotent key).",
	},
	42200: {
		Status: "UnprocessableEntity",
		Message: "The request body was well-formed but contains semantical errors. " +
			"The response body will provide more details in the errors or error parameters.",
	},
	42900: {
		Status:  "TooManyRequests",
		Message: "The request was not accepted because the application has exceeded the rate limit. The default API call limit is 10 requests per second.",
	},
	50000: {
		Status:  "InternalError",
		Message: "Something went wrong on AfterShip's end. Also, some error that cannot be retried happened on an external system that this call relies on.",
	},
}

// HttpCodeDescription returns the description of the http code.
func HttpCodeDescription(code int) Description {
	return code2Description[code]
}

// RegisterCustomCode was used to register the custom code description,
// it would overwrite the previous one when conflicted.
func RegisterCustomCode(code int, status, message string) error {
	if status == "" {
		desc, ok := code2Description[code/100*100]
		if !ok {
			return errors.New("status shouldn't be empty")
		}
		status = desc.Status
	}
	if message == "" {
		return errors.New("message shouldn't be empty")
	}
	code2Description[code] = Description{
		Status:  status,
		Message: message,
	}
	return nil
}
