package nex_utility

import (
	"github.com/PretendoNetwork/mario-kart-7-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/utility"
)

func GetAssociatedNexUniqueIDWithMyPrincipalID(err error, client *nex.Client, callID uint32) {
	// TODO - This method has a different behavior on MK7 compared to the docs.
	// Find out what are the request and response contents
	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)

	rmcResponseStream.WriteUInt64LE(0)
	rmcResponseStream.WriteUInt64LE(0)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(utility.ProtocolID, callID)
	rmcResponse.SetSuccess(utility.MethodGetAssociatedNexUniqueIDWithMyPrincipalID, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(responsePacket)
}
