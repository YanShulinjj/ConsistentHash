/* ----------------------------------
*  @author suyame 2022-06-27 15:37:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package main

import (
	"fmt"
	"github.com/YanShulinJJ/ConsistentHash"
	"log"
	"os"
	"time"
)

type myData struct {
	int
	string
}

func main() {
	hc := ConsistentHash.NewHashCycle()
	// Init log
	l := log.New(os.Stdout, "consistenceHash ", log.Ldate|log.Ltime)
	hc.SetLogger(l)
	// Use heartbeat
	go ConsistentHash.Heartbeat(hc, &map[ConsistentHash.HashAddrType]bool{})
	// Add a new server
	server := ConsistentHash.NewServer("127.0.0.1")
	hc.AddServer(server)
	// Add some items: k-v
	hc.AddItem("key1", "value1")
	hc.AddItem("key2", myData{18, "very good man"})
	// Get value
	v1, err := hc.GetValue("key1")
	if err != nil {
		panic(err)
	}
	fmt.Println(v1)

	v2, err := hc.GetValue("key2")
	if err != nil {
		panic(err)
	}
	fmt.Println(v2)

	// Add another server
	server2 := ConsistentHash.NewServer("192.168.0.1")
	hc.AddServer(server2)
	// Add some items: k-v
	hc.AddItem("key3", "value3")
	hc.AddItem("key4", myData{22, "very great man"})
	// Get value
	v3, err := hc.GetValue("key3")
	if err != nil {
		panic(err)
	}
	fmt.Println(v3)

	v4, err := hc.GetValue("key4")
	if err != nil {
		panic(err)
	}
	fmt.Println(v4)
	// Print internal allocate.
	hc.PrintStatus()
	// server2 crash at this time
	server2.Crash()
	// wait a moment
	time.Sleep(500 * time.Microsecond)
	hc.PrintStatus()
	// server2 restart at this time
	server2.Restart()
	// wait a moment
	time.Sleep(500 * time.Microsecond)
	hc.PrintStatus()
}
