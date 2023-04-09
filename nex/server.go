package nex

import (
	"fmt"

	"github.com/PretendoNetwork/mario-kart-7-secure/database"
	"github.com/PretendoNetwork/mario-kart-7-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
)

func StartNEXServer() {
	globals.NEXServer = nex.NewServer()
	globals.NEXServer.SetPRUDPVersion(0)
	globals.NEXServer.SetDefaultNEXVersion(&nex.NEXVersion{
		Major: 3,
		Minor: 4,
		Patch: 17,
	})
	globals.NEXServer.SetKerberosKeySize(32)
	globals.NEXServer.SetAccessKey(globals.Config.AccessKey)
	globals.NEXServer.SetPingTimeout(999)
	globals.NEXServer.SetKerberosPassword("test")

	globals.NEXServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Println("==MK8 - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Printf("Method ID: %#v\n", globals.NEXServer.NEXVersion())
		fmt.Println("=================")
	})

	globals.NEXServer.On("Kick", func(packet *nex.PacketV0) {
		pid := packet.Sender().PID()
		database.RemovePlayer(pid)

		fmt.Println("Leaving")
	})

	registerCommonProtocols()
	registerNEXProtocols()

	globals.NEXServer.Listen(":" + globals.Config.ServerPort)
}
