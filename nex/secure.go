package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewPRUDPServer()
	globals.SecureServer.SecureVirtualServerPorts = []uint8{1}
	globals.SecureServer.PRUDPVersion = 0
	globals.SecureServer.SetDefaultLibraryVersion(nex.NewLibraryVersion(2, 4, 3))

	globals.SecureServer.SetAccessKey("6181dff1")
	globals.SecureServer.SetKerberosPassword([]byte(globals.KerberosPassword))

	globals.SecureServer.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("=== MK7 - Secure ===")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID)
		fmt.Printf("Method ID: %#v\n", request.MethodID)
		fmt.Println("====================")
	})

	registerCommonSecureServerProtocols()
	registerSecureServerNEXProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_MK7_SECURE_SERVER_PORT"))

	globals.SecureServer.Listen(port)
}
