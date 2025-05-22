package response

import (
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/google/uuid"
)

type Response map[string]any

func Success(data any, message string) Response {
	// msg := "Your request has been successfully processed"
	// // msg := "Permintaan anda berhasil diproses"
	// if message != "" {
	// 	msg = message
	// }

	// if data == nil {
	// 	return Response{
	// 		"success": true,
	// 		"message": msg,
	// 	}
	// }

	return Response{
		// "success": true,
		// "message": msg,
		"data": data,
	}
}

func Error(errorMsg any) Response {
	requestID := uuid.New().String()

	if _, ok := errorMsg.(string); ok {
		return Response{
			"request_id": requestID,
			"errors":     errorMsg,
			// "success": false,
			// "message": errorMsg,
		}
	}

	if _, ok := errorMsg.(map[string][]string); ok {
		return Response{
			// "success": false,
			"request_id": requestID,
			"errors":     errorMsg,
			// "message": "Your request has been failed to process",
			// "message": "Permintaan anda gagal diproses",
		}
	}

	if errHttp, ok := errorMsg.(*err_msg.CustomError); ok {
		return Response{
			"request_id": requestID,
			// "errors":     errHttp.Errors,
			// "success": false,
			"errors": errHttp.Msg,
		}
	}

	if _, ok := errorMsg.(error); ok {
		return Response{
			"request_id": requestID,
			"errors":     make(map[string][]string),
			// "success": false,
			// "message": err.Error(),
		}
	}

	return Response{
		// "success": false,
		// "message": "Your request has been failed to process",
		// "message": "Permintaan anda gagal diproses",
	}
}
