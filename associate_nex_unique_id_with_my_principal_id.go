package main
import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func associateNexUniqueIDsWithMyPrincipalID(err error, client *nex.Client, callID uint32, uniqueIDInfo []*nexproto.UniqueIDInfo) {
	rmcResponse := nex.NewRMCResponse(nexproto.UtilityProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.UtilityMethodAssociateNexUniqueIdsWithMyPrincipalId, nil)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}

