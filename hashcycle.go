/* ----------------------------------
*  @author suyame 2022-06-21 15:25:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package ConsistentHash

import (
	"log"
	"sync"
)

/**

实现一致性哈希

**/

type HashAddrType uint32

type HashCycle struct {
	sync.RWMutex
	serverCount uint32
	servers     map[HashAddrType]*Server
	// history record servers whose has used. access to recover crash server.
	history map[HashAddrType]*Server
	// serverTable records alive servers, must be ordered！
	serverTable []HashAddrType

	// The logger used for this table.
	logger *log.Logger
}

type Item struct {
	key   interface{}
	value interface{}
	// addr is the result of sha256 with key
	addr HashAddrType
}

func NewHashCycle() *HashCycle {
	return &HashCycle{
		history: make(map[HashAddrType]*Server, 0),
		servers: make(map[HashAddrType]*Server, 0),
		logger:  new(log.Logger),
	}
}

// Exist judge server is or not exist.
func (hc *HashCycle) Exist(server_ip string) bool {
	hc.Lock()
	defer hc.Unlock()
	// get ip
	ip, err := transform(server_ip)
	if err != nil {
		panic(err)
	}
	idx := Binsearch(hc.serverTable, ip)
	return idx < len(hc.serverTable) && hc.serverTable[idx] == ip
}

// AddServer adds a new server into the cycle.
func (hc *HashCycle) AddServer(server *Server) error {
	// judge it is already existed.
	if ok := hc.Exist(server.ipstr); ok {
		hc.log("server: ", server.ip, " already exists!")
		return ServerExistError
	}
	// add the server into cycle
	hc.Lock()
	defer hc.Unlock()
	hc.servers[server.ip] = server
	hc.history[server.ip] = server
	// update the table
	// use binsearch !
	idx := Binsearch(hc.serverTable, server.ip)
	// insert the new addr into the table, ensure ordered!
	// need extra room
	c := make([]HashAddrType, idx+1)
	copy(c[:idx], hc.serverTable[:idx])
	c[idx] = server.ip
	// c = append(c, server.ip)
	hc.serverTable = append(c, hc.serverTable[idx:]...)
	// update count
	hc.serverCount++
	// get near server
	// need transport item
	next_idx := (idx + 1) % int(hc.serverCount)
	near_server := hc.servers[hc.serverTable[next_idx]]
	// transport data
	for k, v := range near_server.items {
		ok := false
		if server.ip < near_server.ip {
			if k < server.ip || k >= near_server.ip {
				ok = true
			}
		} else {
			if k < server.ip && k > near_server.ip {
				ok = true
			}
		}
		if ok {
			server.items[k] = v
			delete(near_server.items, k)
		}
	}
	return nil
}

// RemoveServer remove the server which ip is `server_ip`.
func (hc *HashCycle) RemoveServer(server_ip string) (*Server, error) {
	// Judge it does or not exits
	if ok := hc.Exist(server_ip); !ok {
		hc.log("server: ", server_ip, "not exist!")
		return nil, ServerNotFoundError
	}
	// Remove the server
	hc.Lock()
	defer hc.Unlock()
	// get ip
	ip, err := transform(server_ip)
	if err != nil {
		panic(err)
	}
	// get near server ip
	idx := Binsearch(hc.serverTable, ip)
	next_idx := (idx + 1) % int(hc.serverCount)
	server := hc.servers[hc.serverTable[next_idx]]
	//
	drop := hc.servers[hc.serverTable[idx]]
	// transport data
	for k, v := range drop.items {
		server.items[k] = v
	}
	// server.items = append(server.items, drop.items...)

	// remove server
	delete(hc.servers, ip)

	// remove the server of table, ensure ordered!
	hc.serverTable = append(hc.serverTable[:idx], hc.serverTable[idx+1:]...)

	// update count
	hc.serverCount--

	return drop, nil
}

// FindTgtServer  find the server of item.
// This function needs apply rlock first!
// return prams:
//      1. tgtServeraddr: server_ip(uint32)
//      2. tgtItemaddr:   item_hash256_head(uint32)
func (hc *HashCycle) FindTgtServer(datakey string) (tgtServerAddr, tgtItemAddr HashAddrType, err error) {

	if len(hc.serverTable) < 1 {
		return 0, 0, ServerNotFoundError
	}
	tgtServerAddr = hc.serverTable[0]
	tgtItemAddr = GetKeyHashAddr(datakey)
	for _, addr := range hc.serverTable {
		if addr > tgtItemAddr {
			tgtServerAddr = addr
			break
		}
	}
	return tgtServerAddr, tgtItemAddr, nil
}

// AddItem adds Data into system
func (hc *HashCycle) AddItem(key string, value interface{}) error {
	// judge it does or not exist
	// Find the server will save this item.
	hc.RLock()
	tgtServeraddr, tgtItemaddr, err := hc.FindTgtServer(key)
	if err != nil {
		return err
	}
	// check it is or not already exits.
	ok := hc.servers[tgtServeraddr].Exist(tgtItemaddr)
	hc.RUnlock()
	if ok {
		return ItemExistError
	}
	//  Write Item into system
	hc.Lock()
	defer hc.Unlock()
	hc.servers[tgtServeraddr].items[tgtItemaddr] = &Item{
		key:   key,
		value: value,
		addr:  tgtItemaddr,
	}
	// hc.log("Write key: ", key, " value: ", value)
	return nil
}

// GetItem get a item from system.
func (hc *HashCycle) GetItem(key string) (*Item, error) {
	// Find the server. may be better in some way.
	hc.RLock()
	defer hc.RUnlock()
	tgtServeraddr, tgtItemaddr, err := hc.FindTgtServer(key)
	if err != nil {
		return nil, err
	}
	server := hc.servers[tgtServeraddr]
	item, ok := server.items[tgtItemaddr]
	if !ok {
		return nil, ItemNotFoundError
	}

	return item, nil
}

// GetValue get value of key.
func (hc *HashCycle) GetValue(key string) (interface{}, error) {
	item, err := hc.GetItem(key)
	if err != nil {
		return nil, err
	}
	return item.value, nil
}

// SetLogger sets the logger to be used by this cache table.
func (hc *HashCycle) SetLogger(logger *log.Logger) {
	hc.Lock()
	defer hc.Unlock()
	hc.logger = logger
}

// Internal logging method for convenience.
func (hc *HashCycle) log(v ...interface{}) {
	if hc.logger == nil {
		return
	}

	hc.logger.Println(v...)
}

// Get status of hc
func (hc *HashCycle) GetStatus() Status {
	status := Status{
		serverCount: hc.serverCount,
	}
	nodes := make(map[HashAddrType][]HashAddrType)
	for _, serverid := range hc.serverTable {
		server, ok := hc.servers[serverid]
		if !ok {
			panic("ServerTable error!")
		}
		// print out items
		for k, _ := range server.items {
			nodes[serverid] = append(nodes[serverid], k)
		}
	}
	status.nodes = nodes
	return status
}

func (hc *HashCycle) PrintStatus() {
	// print out servers
	hc.log("|****************************** status ********************************")
	hc.log("current_server_num: ", hc.serverCount)
	for _, serverid := range hc.serverTable {
		server, ok := hc.servers[serverid]
		if !ok {
			panic("ServerTable error!")
		}
		hc.log("|--server: ", server.ipstr, "--ip(uint32): ", server.ip)
		// print out items
		for k, v := range server.items {
			hc.log("|----item: ", k, " --key:  ", v.key, " --value: ", v.value)
		}
	}
	hc.log("")
}
