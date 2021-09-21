package packets

const (
	CLIENT_STAR = uint32(0x53544152)
	CLIENT_SSHR = uint32(0x52485353)
	CLIENT_SEXP = uint32(0x50584553)
	CLIENT_JSTR = uint32(0x5254534a)
	CLIENT_W2BN = uint32(0x5732424e)
	CLIENT_WAR3 = uint32(0x33524157)
	CLIENT_W3XP = uint32(0x50583357)
	CLIENT_DSHR = uint32(0x52485344)
	CLIENT_DRTL = uint32(0x4c545244)
	CLIENT_D2DV = uint32(0x56443244)
	CLIENT_D2XP = uint32(0x44325850)
)

var CLIENT_CONFIG = map[uint32]map[string]uint32{
	CLIENT_STAR: {"supported": 0x0, "version": uint32(0xa5)},
	CLIENT_SSHR: {"supported": 0x0, "version": uint32(0xd3)},
	CLIENT_SEXP: {"supported": 0x0, "version": uint32(0xd3)},
	CLIENT_JSTR: {"supported": 0x0, "version": uint32(0xa9)},
	CLIENT_W2BN: {"supported": 0x0, "version": uint32(0x4f)},
	CLIENT_WAR3: {"supported": 0x0, "version": uint32(0x1e)},
	CLIENT_W3XP: {"supported": 0x0, "version": uint32(0x1e)},
	CLIENT_DSHR: {"supported": 0x0, "version": uint32(0xd3)},
	CLIENT_DRTL: {"supported": 0x0, "version": uint32(0x2a)},
	CLIENT_D2DV: {"supported": 0x0, "version": uint32(0x0e)},
	CLIENT_D2XP: {"supported": 0x1, "version": uint32(0x0e)},
}

const PACKET_MARKER = uint8(0xff)

const (
	SID_NULL       = uint8(0x00)
	SID_MESSAGEBOX = uint8(0x19)
	SID_PING       = uint8(0x25)
	SID_AUTH_INFO  = uint8(0x50)
)
