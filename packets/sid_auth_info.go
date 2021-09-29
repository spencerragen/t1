package packets

import (
	"encoding/hex"
	"fmt"
	"net"

	"t1/logging"
	"t1/utils"
)

const LOGON_XSHA1 uint8 = 0x00
const LOGON_NLS1 uint8 = 0x01
const LOGON_NLS2 uint8 = 0x02

// Definition of the SID_AUTH that a client will send to the server
type BNCS_CLIENT_SID_AUTH_INFO struct {
	BNCSBase
	ProtocolID     uint32 // only ever 0x00
	PlatformCode   uint32
	ProductCode    uint32
	Version        uint32
	LanguageCode   uint32
	LocalIP        uint32
	TimeZoneBias   uint32
	MPQLocaleID    uint32
	UserLanguageID uint32
	CountryAbbr    string
	Country        string
}

// Definition of the SID_AUTH packet that the server will respond with
type BNCS_SERVER_SID_AUTH_INFO struct {
	BNCSBase
	LogonType       uint32
	ServerToken     uint32
	UDPValue        uint32
	CR_MPQ_Filetime uint64
	CR_MPQ_Filename string
	CR_Formula      string
	ServerSignature [128]byte // WAR3/W3XP only
}

func (d BNCS_SERVER_SID_AUTH_INFO) String() string {
	return fmt.Sprintf("%x:%x:%x -> %d", d.Marker, d.ID, d.Length, d.Length)
}

func (d *BNCS_SERVER_SID_AUTH_INFO) From(packet BNCSGeneric) {
	// always
	d.Marker = packet.Marker
	d.ID = packet.ID
	d.Length = packet.Length

	// specific
	d.LogonType = packet.ReadUint32()
	d.ServerToken = packet.ReadUint32()
	d.UDPValue = packet.ReadUint32()
	d.CR_MPQ_Filetime = packet.ReadUint64()
	d.CR_MPQ_Filename = packet.ReadString()
	d.CR_Formula = packet.ReadString()

	logging.Infoln("Received SID_AUTH_INFO\n", hex.Dump(utils.GetBytes(packet)))
}

func (d *BNCS_CLIENT_SID_AUTH_INFO) From(p BNCSGeneric) {
	// always
	d.Marker = p.Marker
	d.ID = p.ID
	d.Length = p.Length

	// specific
	d.ProtocolID = p.ReadUint32()
	d.PlatformCode = p.ReadUint32()
	d.ProductCode = p.ReadUint32()
	d.Version = p.ReadUint32()
	d.LanguageCode = p.ReadUint32()
	d.LocalIP = p.ReadUint32()
	d.TimeZoneBias = p.ReadUint32()
	d.MPQLocaleID = p.ReadUint32()
	d.UserLanguageID = p.ReadUint32()
	d.CountryAbbr = p.ReadString()
	d.Country = p.ReadString()
}

// Process a SID_AUTH packet from the client, construct a response, and return it
func (d BNCS_CLIENT_SID_AUTH_INFO) Process(setLocalIp *string) (BNCSGeneric, error) {
	if CLIENT_CONFIG[d.ProductCode]["supported"] == 0x0 {
		return BNCSGeneric{}, fmt.Errorf("game 0x%x configured unsupported", d.ProductCode)
	}

	if d.ProtocolID != DEF_PROTOCOL_ID {
		return BNCSGeneric{}, fmt.Errorf("protocol ID 0x%x is invalid", d.ProtocolID)
	}

	found := false
	for _, v := range DEF_ALLOWED_PLATFORMS {
		if d.PlatformCode == v {
			found = true
			break
		}
	}
	if !found {
		return BNCSGeneric{}, fmt.Errorf("platform 0x%x not permitted", d.PlatformCode)
	}

	if d.Version != CLIENT_CONFIG[d.ProductCode]["version"] {
		return BNCSGeneric{}, fmt.Errorf(
			"version code 0x%x invalid (expected 0x%x)",
			d.Version,
			CLIENT_CONFIG[d.ProductCode]["version"],
		)
	}

	// check the language code here?

	ip_bytes := utils.Uint32ToBytes(d.LocalIP)
	// I with there was variable unpacking in go
	conv := net.IPv4(ip_bytes[0], ip_bytes[1], ip_bytes[2], ip_bytes[3]).String()
	setLocalIp = &conv

	// there should be a check here but I don't feel like figuring out uint32 -> int16 right now
	logging.Println("timezone bias: ", d.TimeZoneBias)

	return BNCSGeneric{}, nil
}
