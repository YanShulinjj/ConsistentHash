/* ----------------------------------
*  @author suyame 2022-06-27 9:28:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package ConsistentHash

import "sync"

type Server struct {
	sync.RWMutex
	ipstr   string
	ip      HashAddrType
	items   map[HashAddrType]*Item
	isAlive bool
}

// NewServer new a server object.
func NewServer(ipstr string) *Server {
	ip, err := transform(ipstr)
	if err != nil {
		panic(err)
	}
	return &Server{
		ipstr:   ipstr,
		ip:      ip,
		items:   make(map[HashAddrType]*Item),
		isAlive: true,
	}
}

// Exist judge if the key exist.
func (s *Server) Exist(key HashAddrType) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.items[key]
	return ok
}

// Crash stimulate server crash.
func (s *Server) Crash() {
	s.Lock()
	defer s.Unlock()
	s.isAlive = false
}

// Restart stimulate crashed server restart.
func (s *Server) Restart() {
	s.Lock()
	defer s.Unlock()
	s.isAlive = true
}
