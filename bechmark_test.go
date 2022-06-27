/* ----------------------------------
*  @author suyame 2022-06-25 15:44:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package ConsistentHash

import (
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkNewHashCycle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// add a new hashcycle
		hc := NewHashCycle()
		// set logger
		l := log.New(os.Stdout, "consistenceHash ", log.Ldate|log.Ltime)
		hc.SetLogger(l)
	}
}

func BenchmarkAddServer(b *testing.B) {
	// add a new hashcycle
	// generate pprof file
	f, _ := os.OpenFile("cpu.pprof", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	hc := NewHashCycle()
	// set logger
	l := log.New(os.Stdout, "consistenceHash ", log.Ldate|log.Ltime)
	hc.SetLogger(l)
	// first add a server into hc
	server := NewServer("127.0.0.1")
	hc.AddServer(server)

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
				b.Error("Error Add item. ", err)
			}
			// check value is correct
			pred, err := hc.GetValue(key)
			if err != nil || pred.(string) != value {
				b.Error("Error Add item. get wrong value or ", err)
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < b.N; i++ {
		// add servers
		// stauts_add recorded status of server load when add new servers.
		n := uint32(i)
		p1 := (n & (255 << 24)) >> 24
		p2 := (n & (255 << 16)) >> 16
		p3 := (n & (255 << 8)) >> 8
		p4 := n & 255
		host := strconv.Itoa(int(p1)) + "." + strconv.Itoa(int(p2)) + "." + strconv.Itoa(int(p3)) + "." + strconv.Itoa(int(p4))
		server := NewServer(host)
		err := hc.AddServer(server)
		// check it is existent.
		hc.RLock()
		_, ok := hc.servers[server.ip]
		hc.RUnlock()
		if err != nil || !ok {
			b.Error("Error add a new server into hash cycle.", err)
		}
	}

}
