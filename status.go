/* ----------------------------------
*  @author suyame 2022-06-27 9:31:00
*  Crazy for Golang !!!
*  IDE: GoLand
*-----------------------------------*/

package ConsistentHash

import "sort"

type Status struct {
	serverCount uint32
	nodes       map[HashAddrType][]HashAddrType
}

// sort
type Interval []HashAddrType

func (a Interval) Len() int           { return len(a) }
func (a Interval) Less(i, j int) bool { return a[i] < a[j] }
func (a Interval) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (src *Status) IsEqual(tgt *Status) bool {
	if src.serverCount != tgt.serverCount {
		return false
	}
	// check every server items.
	for k, v := range src.nodes {
		tgtv, ok := tgt.nodes[k]
		if !ok {
			return false
		}
		sort.Sort(Interval(tgtv))
		sort.Sort(Interval(v))
		//
		for i := range v {
			if v[i] != tgtv[i] {
				return false
			}
		}
	}
	return true
}
