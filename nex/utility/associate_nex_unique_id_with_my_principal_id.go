package nex_utility

import (
	"github.com/PretendoNetwork/mario-kart-7-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-protocols-go/utility"
)

func AssociateNexUniqueIDsWithMyPrincipalID(err error, client *nex.Client, callID uint32, uniqueIDInfo []*utility.UniqueIDInfo) {
	rmcResponse := nex.NewRMCResponse(utility.ProtocolID, callID)
	rmcResponse.SetSuccess(utility.MethodAssociateNexUniqueIDsWithMyPrincipalID, nil)

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
