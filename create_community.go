package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"fmt"
	"encoding/hex"
)

func createCommunity(err error, client *nex.Client, callID uint32, community *nexproto.PersistentGathering, strMessage string) {
	rmcResponseStream := nex.NewStreamOut(nexServer)
	gid := newCommunity(client.PID(), community.CommunityType, community.Password, community.Attribs, community.ApplicationBuffer, community.ParticipationStartDate.Value(), community.ParticipationEndDate.Value())

	rmcResponseStream.WriteUInt32LE(gid)
	rmcResponseBody := rmcResponseStream.Bytes()
	fmt.Println(hex.EncodeToString(rmcResponseBody))

	rmcResponse := nex.NewRMCResponse(nexproto.MatchmakeExtensionProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchmakeExtensionMethodCreateCommunity, rmcResponseBody)

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
