/* ----------------------------------
*  @author suyame 2022-06-22 15:12:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package ConsistentHash

import (
	"crypto/sha256"
)

func GetKeyHashAddr(key string) HashAddrType {

	hash := sha256.New()
	hash.Write([]byte(key))
	sum := hash.Sum(nil)

	var ans HashAddrType
	// only get head 32bit
	for i := 0; i < 4; i++ {
		ans = ans*256 + HashAddrType(sum[i])
	}
	return ans
}
