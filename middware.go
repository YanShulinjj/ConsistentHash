/* ----------------------------------
*  @author suyame 2022-06-24 8:38:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package ConsistentHash

import (
	"strconv"
	"strings"
)

// transform  trans string type to HashAddrType
// ipstr : "XXX.XXX.XXX", using ipv4 type
func transform(ipstr string) (HashAddrType, error) {
	numstrs := strings.Split(ipstr, ".")
	var addr HashAddrType
	for _, numstr := range numstrs {
		num, err := strconv.Atoi(numstr)
		if err != nil {
			return 0, err
		}
		addr = addr*256 + HashAddrType(num)
	}
	return addr, nil
}

// recover transfor HashAddrType to string.
func recover(ipaddr HashAddrType) string {
	n := uint32(ipaddr)
	p1 := (n & (255 << 24)) >> 24
	p2 := (n & (255 << 16)) >> 16
	p3 := (n & (255 << 8)) >> 8
	p4 := n & 255
	host := strconv.Itoa(int(p1)) + "." + strconv.Itoa(int(p2)) + "." + strconv.Itoa(int(p3)) + "." + strconv.Itoa(int(p4))
	return host
}

func Binsearch(arr []HashAddrType, t HashAddrType) int {
	// find the smallest of bigger than t
	l, r := 0, len(arr)
	mid := 0
	for l < r {
		mid = (l + r) >> 1
		if arr[mid] >= t {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}
