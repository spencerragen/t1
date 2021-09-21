package packets

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
)

type BNCS_SERVER_SID_MESSAGEBOX struct {
	BNCSBase
	Style   uint32 // set to 0 for OK only
	Text    string
	Caption string
}

func (d BNCS_SERVER_SID_MESSAGEBOX) String() string {
	return fmt.Sprintf("%x:%x:%x -> %d", d.Marker, d.ID, d.Length, d.Length)
}

// Send SID_MESSAGEBOX to the client. The parameters are passed (in official clients) to
// the Win32 MessageBox API: https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-messagebox?redirectedfrom=MSDN
func SendMessageBox(conn *net.TCPConn, text string, caption string) {
	p := BNCS_SERVER_SID_MESSAGEBOX{}
	p.Marker = PACKET_MARKER
	p.ID = SID_MESSAGEBOX
	p.Length = uint16(10 + len(text) + len(caption))
	p.Style = uint32(0x00)
	p.Text = text
	p.Caption = caption

	log.Println("Sending SID_MESSAGEBOX\n", hex.Dump(GetBytes(p)))

	_, err := conn.Write(GetBytes(p))
	if err != nil {
		log.Println("[!] failed to write SID_MESSAGEBOX to client", conn.RemoteAddr().String(), ":", err.Error())
		conn.Close()
	}
}
