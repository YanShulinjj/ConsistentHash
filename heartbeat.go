/* ----------------------------------
*  @author suyame 2022-06-27 9:34:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package ConsistentHash

import (
	"time"
)

// Heartbeat
func Heartbeat(hc *HashCycle, sa *map[HashAddrType]bool) error {
	// get the servers alive value now.
	sa_after := make(map[HashAddrType]bool)
	for k, v := range hc.history {
		sa_after[k] = v.isAlive
		// if changed.
		// 1. added new server, do not anything.
		if live, ok := (*sa)[k]; !ok {
			continue
		} else if live != v.isAlive {
			//
			hc.log("[Tick] servers status changed!")
			if v.isAlive {
				// restarted a server
				// we need recover it
				hc.AddServer(v)
			} else {
				// the server crash!
				// we need remove it
				hc.RemoveServer(v.ipstr)
			}

		}
	}
	// ervery 100ms read the servers status.
	time.AfterFunc(100*time.Microsecond, func() {
		Heartbeat(hc, &sa_after)
	})
	return nil
}
