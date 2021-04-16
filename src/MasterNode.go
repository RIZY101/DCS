package main

//@Author: Richard Zins

import (
	"fmt"
	"crypto/rand"
	"crypto/tls"
	"math/big"
	"os/exec"
	"net"
	"log"
	"strings"
	"strconv"
	"sync"
)

//Each nodeId is mapped to one of these structs
//Also storageAvailable is in GB
type NodeData struct {
	ip string
	key string
	storageAvailable float64
}

//Global variables are bad dont use them...
var mapOfNodes map[string]NodeData = make(map[string]NodeData)
var mutex sync.Mutex

func main() {

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
	//6633 is node on the keypad 	
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
    ip := c.RemoteAddr().String()
	log.Printf("Client connected! Their addr is: " + ip)
	buffer := make([]byte, 64)
	c.Read(buffer)
	msg := string(buffer)
	mutex.Lock()
	resp := parseMsg(msg, ip)
	mutex.Unlock()
	c.Write([]byte(resp))
}

func parseMsg(msg string, ip string) string {
	args := strings.Split(msg, " ")

	if args[0] == "ATL" && len(args) == 2 {
		nodeId := strings.Trim(genNodeId(), " ")
		nodeId = strings.Trim(nodeId, "\x0a")
		key := genKey()
		storStr := strings.Trim(args[1], "\x00")
		storage, err  := strconv.ParseFloat(storStr, 64)
		if err != nil {
			log.Fatal(err)
		}
		mapOfNodes[nodeId] = NodeData{ip, key, storage}
		log.Printf("VALID REQUEST")
		//Right now this will make a new nodeId and always allow them to join even if the IP is already another nodeID
		return "ATLR yes " + nodeId + " " + key
	} else if args[0] == "RFL" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		if mapOfNodes[args[1]].ip != "" && strings.Trim(args[2], "\x00") == mapOfNodes[args[1]].key {
			delete(mapOfNodes, args[1])
			return "RFLR yes"
		}
		return "RFLR No"
	} else if args[0] == "NODE" && len(args) == 1 {
		log.Printf("VALID REQUEST")
		return "NODER ipOfNewNode nodeId\n"
	} else if args[0] == "UPDATE" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		return "UPDATER yesOrNO"
	} else {
		log.Printf("NOT A VALID REQUEST")
	}
	return "NOT A VALID REQUEST"
}

//Please note this function uses a unix utility to generate a standard uuid
//You can learn more about uuid's here: https://en.wikipedia.org/wiki/Universally_unique_identifier
func genNodeId() string {
	nodeId, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(nodeId)
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

func printMap() {
	//Iterates over all keys and values in a map
	for k, v := range mapOfNodes {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
}


