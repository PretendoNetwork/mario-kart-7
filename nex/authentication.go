package nex

import (
	"fmt"
	"os"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/mario-kart-7/globals"
)

var serverBuildString string

func StartAuthenticationServer() {
	globals.AuthenticationServer = nex.NewPRUDPServer()
	globals.AuthenticationServer.PRUDPVersion = 0
	globals.AuthenticationServer.SetDefaultLibraryVersion(nex.NewLibraryVersion(2, 4, 3))

	globals.AuthenticationServer.SetKerberosPassword([]byte(globals.KerberosPassword))
	globals.AuthenticationServer.SetAccessKey("6181dff1")

	globals.AuthenticationServer.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("=== MK7 - Auth ===")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID)
		fmt.Printf("Method ID: %#v\n", request.MethodID)
		fmt.Println("==================")
	})

	registerCommonAuthenticationServerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_MK7_AUTHENTICATION_SERVER_PORT"))

	globals.AuthenticationServer.Listen(port)
}
