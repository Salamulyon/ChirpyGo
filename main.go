package main

import (
	"net/http"
)





func main (){

	serveMux := http.serveMux()
	server := Server{
		handler : serveMux,
		Addr : ":8080"
	}
}