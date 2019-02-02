package functions

import (
	"context"
	"net/http"
)

func Response(res http.ResponseWriter, req *http.Request, status int, json []byte) {
	*req = *req.WithContext(context.WithValue(req.Context(), "statusCode", status))
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(status)
	res.Write(json)
	return
}
