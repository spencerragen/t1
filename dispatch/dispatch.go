package dispatch

import (
	"fmt"
	"t1/packets"
)

const (
	SID_AUTH_INFO uint8 = 0x50
)

func Dispatch(packet packets.BNCSGeneric) error {
	switch packet.ID {
	case SID_AUTH_INFO:
		fmt.Println("received SID_AUTH_INFO")

	}
	return nil
}
