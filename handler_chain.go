package main

import "net/http"

type HandlerChain []http.Handler

func (chain HandlerChain) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, v := range chain {
		v.ServeHTTP(writer, request)
	}
}
