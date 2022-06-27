/* ----------------------------------
*  @author suyame 2022-06-24 14:36:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package ConsistentHash

import (
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	// add a new hashcycle
	hc := NewHashCycle()
	// set logger
	l := log.New(os.Stdout, "consistenceHash ", log.Ldate|log.Ltime)
	hc.SetLogger(l)
	wg := sync.WaitGroup{}
	N := 255 // this number should less than 256.
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(j int) {
			// add a server into hc
			// new a server
			defer wg.Done()
			hostsquence := strconv.Itoa(j)
			server := NewServer("192.168.0." + hostsquence)
			err := hc.AddServer(server)
			// check it is existent.
			hc.RLock()
			_, ok := hc.servers[server.ip]
			hc.RUnlock()
			if err != nil || !ok {
				t.Error("Error add a new server into hash cycle.", err)
			}
			// check hc.recordTable
			for _, addr := range hc.serverTable {
				if addr == server.ip {
					goto OUTBREAK
					break
				}
			}
			t.Error("Error add a new server into hash cycle. table hasn't record!")
		OUTBREAK:
			log.Println("Sucessfully add server: ", server.ipstr)
		}(i)
	}
	wg.Wait()
	// check sum of servers.
	if len(hc.serverTable) != N {
		t.Error("Error add server, sum of servers doesn't matched! ")
	}
}

func TestFindTgtServer(t *testing.T) {
	// add a new hashcycle
	hc := NewHashCycle()
	// set logger
	l := log.New(os.Stdout, "consistenceHash ", log.Ldate|log.Ltime)
	hc.SetLogger(l)
	// first add a server into hc
	server := NewServer("192.168.0.1")
	hc.AddServer(server)
	key1 := "suyame"
	// get key1
	server_ip, _, err := hc.FindTgtServer(key1)
	if err != nil || server_ip != server.ip {
		t.Error("Error find target server. ", err)
	}
	// add a new server into hc
	server = NewServer("192.168.0.2")
	hc.AddServer(server)
	key2 := "junjun"
	server_ip, key_addr, err := hc.FindTgtServer(key2)
	if err != nil {
		t.Error("Error find target server. ", err)
	}
	// judge key_addr is biggest of those less than server_ip
	var correct HashAddrType = hc.serverTable[0]
	for _, addr := range hc.serverTable {
		if addr > key_addr {
			correct = addr
			break
		}
	}
	if correct != server_ip {
		t.Error("Error find target server, find server is not matched!", key_addr, correct, server_ip)
	}
}

func TestAddItem(t *testing.T) {
	// add a new hashcycle
	hc := NewHashCycle()
	// set logger
	l := log.New(os.Stdout, "consistenceHash ", log.Ldate|log.Ltime)
	hc.SetLogger(l)
	// first add a server into hc
	server := NewServer("192.168.0.1")
	hc.AddServer(server)

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		// add a key- value item
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(j)
			value := "value" + strconv.Itoa(j)
			err := hc.AddItem(key, value)
			if err != nil {
				t.Error("Error Add item. ", err)
			}
			// check value is correct
			pred, err := hc.GetValue(key)
			if err != nil || pred.(string) != value {
				t.Error("Error Add item. get wrong value or ", err)
			}
		}(i)
	}
	wg.Wait()
	// get values
	for i := 0; i < 100; i++ {
		key := "key" + strconv.Itoa(i)
		_, err := hc.GetValue(key)
		if err != nil {
			t.Error("Error Add item. get wrong value or ", err)
		}
	}
}

// TestRemoveServer check if it does work when some server crash!
func TestRemoveServer(t *testing.T) {
	// add a new hashcycle
	hc := NewHashCycle()
	// set logger
	l := log.New(os.Stdout, "consistenceHash ", log.Ldate|log.Ltime)
	hc.SetLogger(l)
	// first add tow new server into hc
	server1 := NewServer("192.168.0.1")
	hc.AddServer(server1)
	server2 := NewServer("192.168.0.2")
	hc.AddServer(server2)
	// add some key-value items.
	N := 100
	for i := 0; i < N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		err := hc.AddItem(key, value)
		if err != nil {
			t.Error("Error add k-v error, ", err)
		}
	}
	// server1 crash !
	_, err := hc.RemoveServer(server1.ipstr)
	if err != nil {
		t.Error("Error remove server, ", err)
	}
	// check k-v
	for i := 0; i < N; i++ {
		key := "key" + strconv.Itoa(i)
		value := "value" + strconv.Itoa(i)
		pred, err := hc.GetValue(key)
		if err != nil || pred.(string) != value {
			t.Error("Error can't get value when server1 crash, ", err)
		}
	}
}

// TestAddServer  check if it does work when some server restarts!
func TestAddServer(t *testing.T) {
	// add a new hashcycle
	hc := NewHashCycle()
	// set logger
	l := log.New(os.Stdout, "consistenceHash ", log.Ldate|log.Ltime)
	hc.SetLogger(l)
	// first add a server into hc
	server := NewServer("127.0.0.1")
	hc.AddServer(server)

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		// add a key- value item
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(j)
			value := "value" + strconv.Itoa(j)
			err := hc.AddItem(key, value)
			if err != nil {
				t.Error("Error Add item. ", err)
			}
			// check value is correct
			pred, err := hc.GetValue(key)
			if err != nil || pred.(string) != value {
				t.Error("Error Add item. get wrong value or ", err)
			}
		}(i)
	}
	wg.Wait()
	// add servers
	N := 255 // this number should less than 256.
	// stauts_add recorded status of server load when add new servers.
	stauts_add := []Status{}
	for i := 0; i < N; i++ {
		hostsquence := strconv.Itoa(i)
		server := NewServer("192.168.0." + hostsquence)
		err := hc.AddServer(server)
		// check it is existent.
		hc.RLock()
		_, ok := hc.servers[server.ip]
		hc.RUnlock()
		if err != nil || !ok {
			t.Error("Error add a new server into hash cycle.", err)
		}
		stauts_add = append(stauts_add, hc.GetStatus())
	}
	// stauts_remove recorded status of server load when remove  servers.
	stauts_remove := []Status{}
	for i := N - 1; i >= 0; i-- {
		stauts_remove = append(stauts_remove, hc.GetStatus())
		hostsquence := strconv.Itoa(i)
		_, err := hc.RemoveServer("192.168.0." + hostsquence)
		if err != nil {
			t.Error("Error remove server, ", err)
		}
	}
	// check status is matched
	n := len(stauts_remove)
	for i, statue := range stauts_add {
		if !statue.IsEqual(&stauts_remove[n-i-1]) {
			t.Error("Error status check failed!")
		}
	}
}

func TestHeartbeat(t *testing.T) {
	ip := "192.168.110.1"
	hc := NewHashCycle()
	// set logger
	l := log.New(os.Stdout, "consistenceHash ", log.Ldate|log.Ltime)
	hc.SetLogger(l)
	// add server
	server1 := NewServer(ip)
	go Heartbeat(hc, &map[HashAddrType]bool{})
	//
	hc.AddServer(server1)
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		// add a key- value item
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(j)
			value := "value" + strconv.Itoa(j)
			err := hc.AddItem(key, value)
			if err != nil {
				log.Panicln("Error Add item. ", err)
			}
			// check value is correct
			pred, err := hc.GetValue(key)
			if err != nil || pred.(string) != value {
				log.Panicln("Error Add item. get wrong value or ", err)
			}
		}(i)
	}
	wg.Wait()
	server2 := NewServer("127.0.0.1")
	hc.AddServer(server2)
	// server1 crush
	s1 := hc.GetStatus()
	server1.Crash()
	// wait 1s
	time.Sleep(time.Second)
	//
	s2 := hc.GetStatus()
	// server1 restart
	server1.Restart()
	// wait 1s
	time.Sleep(time.Second)
	s3 := hc.GetStatus()

	if s1.IsEqual(&s2) {
		t.Error("Error remove crashed server.")
	}
	if !s1.IsEqual(&s3) {
		t.Error("Error recover crashed server.")
	}
}
