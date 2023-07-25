package shared

import "github.com/lightningnetwork/lnd/lnwire"

func EmptyMsatIfNil(input *lnwire.MilliSatoshi) int {
	if input == nil {
		return 0
	}
	return int(*input)
}
