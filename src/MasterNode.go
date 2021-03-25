package main

import (
	"crypto/rand"
	"crypto/tls"
	"net"
	"log"
)

func main() {
	//6633 is node on the keypad
	//connThread, err := net.Listen("tcp", "localhost:6633")
	 //if err != nil {
	 	//fmt.Println("Error listening:", err.Error())
	 	//os.Exit(1)
	 //}
	 
	 //defer connThread.Close()

	 //for {
	 	//connection, err := connThread.Accept()
	 	//if err != nil {
	 		//fmt.Println("Error connecting:", err.Error())
	 		//return
	 	//}
	 	//fmt.Println("Client connected! Their addr is: " + connection.RemoteAddr().String())
	//}

	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")

	if err != nil {
		log.Fatal(err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAnyClientCert}
	config.Rand = rand.Reader
     
	listen, err := tls.Listen("tcp", "localhost:6633", &config)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Printf("Server(TLS) up and listening on port 6633")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	log.Printf("Client connected! Their addr is: " + c.RemoteAddr().String())
}


