package main

import (
	"fmt"
	"os"
	"net"
)

func main() {
	//6633 is node on the keypad
	connThread, err := net.Listen("tcp", "localhost:6633")
	 if err != nil {
	 	fmt.Println("Error listening:", err.Error())
	 	os.Exit(1)
	 }
	 
	 defer connThread.Close()

	 for {
	 	connection, err := connThread.Accept()
	 	if err != nil {
	 		fmt.Println("Error connecting:", err.Error())
	 		return
	 	}
	 	fmt.Println("Client connected! Their addr is: " + connection.RemoteAddr().String())
	 }
}
