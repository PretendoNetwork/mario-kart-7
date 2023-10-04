package nex

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go"
)

func StartSecureServer() {
	globals.SecureServer = nex.NewServer()
	globals.SecureServer.SetPRUDPVersion(0)
	globals.SecureServer.SetDefaultNEXVersion(nex.NewNEXVersion(2, 4, 3))

	globals.SecureServer.SetAccessKey("6181dff1")
	globals.SecureServer.SetKerberosPassword(globals.KerberosPassword)

	globals.SecureServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Println("=== MK7 - Secure ===")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("====================")
	})

	registerCommonSecureServerProtocols()
	registerSecureServerNEXProtocols()

	globals.SecureServer.Listen(fmt.Sprintf(":%s", os.Getenv("PN_MK7_SECURE_SERVER_PORT")))
}
