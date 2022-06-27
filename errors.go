/* ----------------------------------
*  @author suyame 2022-06-21 15:53:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package ConsistentHash

import "errors"

var (
	ServerExistError    = errors.New("Server has already existed!")
	ServerNotFoundError = errors.New("Server Not Found!")
	ServerCrashError    = errors.New("Server has crashed!")
	ItemExistError      = errors.New("Item has aleady existed!")
	ItemNotFoundError   = errors.New("Item Not Found!")
)
