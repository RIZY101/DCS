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
     
	listen, err := tls.Listen("tcp", "192.168.254.87:6633", &config)
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
	buffer := make([]byte, 128)
	c.Read(buffer)
	msg := string(buffer)
	log.Printf(msg)
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
		str := strings.Split(ip, ":")
		ipStr := str[0]
		mapOfNodes[nodeId] = NodeData{ipStr, key, storage}
		log.Printf("VALID REQUEST")
		printMap()
		//Right now this will make a new nodeId and always allow them to join even if the IP is already another nodeID
		return "ATLR yes " + nodeId + " " + key
	} else if args[0] == "RFL" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		if mapOfNodes[args[1]].ip != "" && strings.Trim(args[2], "\x00") == mapOfNodes[args[1]].key {
			delete(mapOfNodes, args[1])
			return "RFLR yes"
		}
		printMap()
		return "RFLR no"
	} else if args[0] == "NODE" && len(args) == 2 {
		storStr := strings.Trim(args[1], "\x00")
		storageNeeded, err := strconv.ParseFloat(storStr, 64)
		if err != nil {
			log.Fatal(err)
		}
		ipOfNode, nodeId := getRandomNode(storageNeeded)
		//Here I decrease the amount of avaialble storage that node has
		newData := NodeData{ipOfNode, nodeId, mapOfNodes[nodeId].storageAvailable - storageNeeded}
		mapOfNodes[nodeId] = newData
		log.Printf("VALID REQUEST")
		printMap()
		return "NODER " + ipOfNode + " " + nodeId
	} else if args[0] == "UPDATE" && len(args) == 3 {
		log.Printf("VALID REQUEST")
		if ip != mapOfNodes[args[1]].ip && args[2] == mapOfNodes[args[1]].key {
			newData := NodeData{ip, mapOfNodes[args[1]].key, mapOfNodes[args[1]].storageAvailable}
			mapOfNodes[args[1]] = newData
			printMap()
			return "UPDATER yes"
		}
		printMap()
		return "UPDATER no"
	} else if args[0] == "CHECK" && len(args) == 3 {
			log.Printf("VALID REQUEST")
			if strings.Trim(args[2], "\x00") != mapOfNodes[args[1]].ip {
				printMap()
				return "CHECKR yes " + mapOfNodes[args[1]].ip
			}
			printMap()
			return "CHECKR no 0.0.0.0"
	} else {
		log.Printf("NOT A VALID REQUEST")
		printMap()
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

//TODO Adapt this function so that it cant return itself but I must do this after I set up testing multiple clients at once
//Only use this is in a mutex safe place
func getRandomNode(storage float64) (string, string) {
	var key string
	var node NodeData
	notFound := true
	for notFound {
		key, node = getNode()
		if node.storageAvailable > storage {
			notFound = false
		}
	}
    return node.ip, key
}

////Only use this is in a mutex safe place
func getNode() (string, NodeData) {
	var key string
	var node NodeData
	//This is only a pseduo random solution
	for k, v := range mapOfNodes {
		node = v
		key = k
		break
	}
	return key, node
}

//Only use this is in a mutex safe place
func printMap() {
	//Iterates over all keys and values in a map
	for k, v := range mapOfNodes {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
}
