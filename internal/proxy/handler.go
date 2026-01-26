package proxy

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"strings"
)

func handleConnect(w http.ResponseWriter, r *http.Request) {
	// CONNECT Hijacking
	hijacker, ok := w.(http.Hijacker) // this is taking full control of this connection. From here, Golang will stop managing the connection.
	if !ok {
		log.Println("[+] ERROR :: Hijacker not supported")
		http.Error(w, "Hijacker is not supported", http.StatusInternalServerError)
	}

	// client connection
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		log.Println("[+] ERROR :: Hijacking failed: ", err)
	}

	// connect to destination
	targetConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		log.Println("[+] ERROR :: Target dial has been failed: ", err)
		clientConn.Close() // closing connection
		return
	}

	// VERY IMPORTANT: write directly to hijacked conn
	/*
		Here, we need to manually send message because we are managing a raw TCP connection. Golang has stopped handling this connection.
	*/
	_, err = clientConn.Write([]byte(
		"HTTP/1.1 200 Connection Established\r\n\r\n",
	))
	if err != nil {
		log.Println("[+] ERROR :: Write failed:", err)
		clientConn.Close()
		targetConn.Close()
		return
	}

	// pipe data
	/*
		This is bidirectional piping. I am not handling HTTPS, I am just relying the bytes.
		Client <===> Proxy <===> Target

		go io.Copy(targetConn, clientConn)
		go io.Copy(clientConn, targetConn)
	*/

	// TLS is stuck at this point
	// start with CERT
	// loading cert
	cert, err := tls.LoadX509KeyPair("../../cert/proxy.crt", "../../cert/proxy.key")
	if err != nil {
		log.Println("[+] ERROR :: Cert loading failed :", err)
		clientConn.Close()
		targetConn.Close()
		return
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	tlsClientConn := tls.Server(clientConn, tlsConfig)

	// Client <===> Proxy handshake
	err = tlsClientConn.Handshake()
	if err != nil {
		log.Println("[+] ERROR :: Client Handshake failed :", err)
		tlsClientConn.Close()
		return
	}

	// reading TLS connection data
	reader := bufio.NewReader(tlsClientConn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Println("[+] ERROR :: Failed to read HTTP request :", err)
		tlsClientConn.Close()
		return
	}

	log.Println("[+] TLS handshake with client successful")
	log.Println("[+] Decrypted HTTP Data: ", req.URL.String(), req.Host, req.Method)

	// Proxy <===> TLS <===> Server handshake
	// dial to main target
	serverConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		log.Println("[+] ERROR :: Failed to connect to the target :", err)
		tlsClientConn.Close()
		return
	}

	tlsServerConn := tls.Client(serverConn, &tls.Config{
		ServerName: strings.Split(r.Host, ":")[0],
	})

	err = tlsServerConn.Handshake()
	if err != nil {
		log.Println("[+] Retrying with www. prefix")

		// close connection
		serverConn.Close()

		// re-open connection
		serverConn, err := net.Dial("tcp", r.Host)
		if err != nil {
			log.Println("[+] ERROR :: Failed to connect to the target :", err)
			tlsClientConn.Close()
			return
		}

		tlsServerConn := tls.Client(serverConn, &tls.Config{
			ServerName: "www." + strings.Split(r.Host, ":")[0],
		})

		err = tlsServerConn.Handshake()
		if err != nil {
			log.Println("[+] ERROR :: Server handshake failed", err)
			tlsServerConn.Close()
			return
		}
	}

	log.Println("[+] TLS handshake with server successful")

	// forwarding request to the server
	req.RequestURI = "" // this field is only for proxy, so reset it before sending to the target

	err = req.Write(tlsServerConn)
	if err != nil {
		log.Println("[+] ERROR :: Request forwarding has been failed", err)
		tlsServerConn.Close()
		return
	}

	log.Println("[+] Request has been forwarded to the server")

	// read response and forward to client
	serverReader := bufio.NewReader(tlsServerConn)

	resp, err := http.ReadResponse(serverReader, req)
	if err != nil {
		log.Println("[+] ERROR :: Reading HTPP response has been failed", err)
		tlsServerConn.Close()
		return
	}

	log.Println("[+] Response received from server:", resp.Status)

	err = resp.Write(tlsClientConn)
	if err != nil {
		log.Println("[+] ERROR :: Failed to send response to the client", err)
		tlsServerConn.Close()
		tlsClientConn.Close()
		return
	}

	log.Println("[+] Response sent back to client")
}

// HandleIncomingReq is the entry point for all request.
func HandleIncomingReq(w http.ResponseWriter, r *http.Request) {
	log.Println("======================= INCOMING REQUEST =======================")

	// handling CONNECT
	/*
		What happend after CONNECT?
		Client send CONNECT method through HTTP and tells to open a raw tunnel to example.com. Afetr CONNECT, the all work of HTTP is over.
		Golang's built-in functions like http.ResponseWriter, WriterHeader, ServerMux can't handle this. From here we get into the raw TCP socket connection, Golang will not manage
		from here.
	*/
	if r.Method == http.MethodConnect {
		log.Println("[+] CONNECT received!", r.Host)
		handleConnect(w, r)
		return
	}

	log.Printf("Method: %s\n", r.Method)
	log.Printf("URL: %s\n", r.URL.String())

	w.WriteHeader(http.StatusOK)
	log.Println("Incoming request has been recieved and logged succcessfully by proxy")

	log.Println("================================================================")
}
