package main

//@Author: Richard Zins

import (
	"crypto/tls"
	"log"
	"strings"
	"os"
	"fmt"
	"crypto/rand"
	"strconv"
	"math/big"
	"io/ioutil"
	"net"
)

var NodeId string
//Your key to give to the LNode storing your file
var Key string


type Node struct {
	ip string
	NodeId string
}
//Global variables are bad dont use them...
//This is where your data is stored
//Each file name is the key
var mapOfYourData map[string]Node = make(map[string]Node)


func main() {
	NodeId = ""
	Key = genKey()
	args := os.Args
	//6633 is node on the keypad
    cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
    if err != nil {
    	log.Fatal(err)
    }

	ip := "localhost"
	//For MNode
	//port := "6633"
	//For LNode
	port := "6634"

	log.Printf("Connecting to %s\n", ip)

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", ip+":"+port, &config)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connection established between %s and localhost.\n", conn.RemoteAddr().String()) 

	//conn.Write([]byte("ATL ipOfNode storageInGB"))
	conn.Write([]byte(args[1]))
	buffer := make([]byte, 64)
	conn.Read(buffer)
	log.Printf(string(buffer))
	size := parseMsg(string(buffer), conn)
	if size > 0 {
		buffer2 := make([]byte, size)
		conn.Read(buffer2)
		data := string(buffer2)
		log.Printf(data)
		pwd, _ := os.Getwd()
		f, err := os.Create(pwd + "/data2/test")
		if err != nil {
			log.Fatal(err)
		}
		num, err := f.Write(buffer2)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Wrote %d bytes", num)
	}
	defer conn.Close()
	log.Printf("Connection Killed")    
}

func parseMsg (msg string, conn net.Conn) int {
	args := strings.Split(msg, " ")
		//Send back STORER yesOrNo
	if args[0] == "NODER" && len(args) == 3 {
		//TODO implment adding correct file name later
		str := strings.Split(args[1], ":")
		ip := str[0]
		mapOfYourData["TestFileName"] = Node{ip, args[2]}
		log.Printf("VALID RESPONSE")
		printMap()
	} else if args[0] == "CHECKR" && len(args) == 3 {
		//TODO Implement correct fileName and NodeId later
		if args[1] == "no" {
			log.Printf("No change in the IP address")
		} else {
			str := strings.Split(args[2], ":")
			ip := str[0]
			mapOfYourData["TestFileName"] = Node{ip, "TestNodeId"}
			log.Printf("The IP address changed")
		}
		printMap()
	} else if args[0] == "STORER" && len(args) == 2 {
		//TODO Implement this so it does specific file not test every time
		log.Printf("VALID RESPONSE")
		str := strings.Trim(args[1], "\x00")
		if str == "yes" {
			log.Printf("LNode ready for the data")
		} else {
			log.Printf("LNode not ready for your data")
		}
		data, err := ioutil.ReadFile("test")
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte(data))
	} else if args[0] == "RETRIEVER" && len(args) == 3 {
		//TODO if no tell the user maybe their key was invalid
		log.Printf("VALID RESPONSE")
		sizeF := strings.Trim(args[2], "\x00")
		size, err := strconv.Atoi(sizeF)
		if err != nil {
			log.Fatal(err)
		}
		return size
	} else if args[0] == "REMOVER" && len(args) == 2 {
		//TODO If no tell user maybe key was wrong
		log.Printf("VALID RESPONSE")
		str := strings.Trim(args[1], "\x00")
		if str == "yes" {
			log.Printf("Data Removed")
		} else {
			log.Printf("Failed to remove your data. Maybe the key you sent was wrong.\n")
		}
	} else {
		log.Printf("NOT A VALID REQUEST OR RESPONSE")
	}
	return 0
}

//Usef for testing only
func printGlobals () {
	log.Printf(NodeId + "\n")
	log.Printf(Key + "\n")
}

//Only use this is in a mutex safe place
func printMap() {
	//Iterates over all keys and values in a map
	for k, v := range mapOfYourData {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
}

//crypto/rand implements a cryptographically secure random number generator so we can just use that 
//Also the key should be aprox 64 bits (This might need to get bigger later)
func genKey() string {
	key, err := rand.Int(rand.Reader, big.NewInt(1000000000000000000))
	if err != nil {
		log.Fatal(err)
	}
	keyStr := key.String()
	return keyStr
}

//TODO Remember to send a STORE after recieving a NODER
