package main

import (
	"fmt"
	"os"
	"net"
)

func main() {
	//6633 is node on the keypad
	connThread, err := net.Dial("tcp", "localhost:6633")
	 if err != nil {
	 	fmt.Println("Error connecting:", err.Error())
	 	os.Exit(1)
	 }

	 defer connThread.Close()
}
