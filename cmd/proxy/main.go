package main

import (
	proxy "ProxyServer/internal/proxy"
	"log"
	"net/http"
)

func main() {
	/*
		Previous code -

		http.HandleFunc("/", proxy.HandleIncomingReq)

		port := ":8000"
		log.Printf("Server is listening on %s...", port)

		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatal("[+] ERROR :: Server starting has been failed :: ", err)
		}

		Why this is incompatible for the proxy?
	*/

	port := ":8000"

	server := &http.Server{
		Addr:    port,
		Handler: http.HandlerFunc(proxy.HandleIncomingReq),
	}

	log.Printf("Proxy is listening on %s...", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("[+] ERROR :: Proxy starting has been failed :: ", err)
	}
}
