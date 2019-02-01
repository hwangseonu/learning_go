package common

import (
	"log"
	"net/http"
)

type Middleware func(next http.Handler) http.HandlerFunc

func ChainBeforeMiddlewares(middlewares ...Middleware) Middleware {
	return func(final http.Handler) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			var last http.HandlerFunc = final.ServeHTTP
			for i := len(middlewares) - 1; i >= 0; i-- {
				last = middlewares[i](last)
			}
			last.ServeHTTP(res, req)
		}
	}
}

func ChainAfterMiddlewares(middlewares ...Middleware) Middleware {
	return func(final http.Handler) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			final.ServeHTTP(res, req)
			var last http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {}
			for i := len(middlewares) - 1; i >= 0; i-- {
				last = middlewares[i](last)
			}
			last.ServeHTTP(res, req)
		}
	}
}

func LoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("[%s] %s %d", req.Method, req.URL.Path, req.Context().Value("statusCode"))
		next.ServeHTTP(res, req)
	}
}