package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"fmt"
	"encoding/hex"
)

func getSimpleCommunity(err error, client *nex.Client, callID uint32, gatheringIdList []uint32) {
	rmcResponseStream := nex.NewStreamOut(nexServer)
	lstSimpleCommunityList := make([]*nexproto.SimpleCommunity, 0)
	for _, gid := range gatheringIdList {
		fmt.Println("gid1", gid)
		if(communityExists(gid)){
			fmt.Println("gid2", gid)
			_, _, _, _, _, _, _, sessions, _ := getCommunityInfo(gid)
			simpleCommunity := nexproto.NewSimpleCommunity()
			simpleCommunity.GatheringID = gid
			simpleCommunity.MatchmakeSessionCount = sessions
			lstSimpleCommunityList = append(lstSimpleCommunityList, simpleCommunity)
		}
	}

	rmcResponseStream.WriteUInt32LE(uint32(len(lstSimpleCommunityList)))
	for _, simpleCommunity := range lstSimpleCommunityList {
		rmcResponseStream.WriteStructure(simpleCommunity)
	}
	rmcResponseBody := rmcResponseStream.Bytes()
	fmt.Println(hex.EncodeToString(rmcResponseBody))

	rmcResponse := nex.NewRMCResponse(nexproto.MatchmakeExtensionProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchmakeExtensionMethodGetSimpleCommunity, rmcResponseBody)

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
