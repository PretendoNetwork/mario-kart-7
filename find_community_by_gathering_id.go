package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"fmt"
	"encoding/hex"
)

func findCommunityByGatheringId(err error, client *nex.Client, callID uint32, lstGid []uint32) {
	rmcResponseStream := nex.NewStreamOut(nexServer)
	lstCommunity := make([]*nexproto.PersistentGathering, 0)
	for _, gid := range lstGid {
		fmt.Println("gid1", gid)
		if(communityExists(gid)){
			fmt.Println("gid2", gid)
			host, communityType, password, attribs, application_buffer, start_date, end_date, sessions, participationCount := getCommunityInfo(gid)
			persistentGathering := nexproto.NewPersistentGathering()
			persistentGathering.Gathering.ID = gid
			persistentGathering.Gathering.OwnerPID = host
			persistentGathering.Gathering.HostPID = host
			persistentGathering.Gathering.MinimumParticipants = 1
			persistentGathering.Gathering.MaximumParticipants = 8
			persistentGathering.CommunityType = communityType
			persistentGathering.Password = password
			persistentGathering.Attribs = attribs
			persistentGathering.ApplicationBuffer = application_buffer
			persistentGathering.ParticipationStartDate = nex.NewDateTime(start_date)
			persistentGathering.ParticipationEndDate = nex.NewDateTime(end_date)
			persistentGathering.ApplicationBuffer = application_buffer
			persistentGathering.MatchmakeSessionCount = sessions
			persistentGathering.ParticipationCount = participationCount
			lstCommunity = append(lstCommunity, persistentGathering)
		}
	}

	rmcResponseStream.WriteUInt32LE(uint32(len(lstCommunity)))
	fmt.Println("len", len(lstCommunity))
	for _, persistentGathering := range lstCommunity {
		rmcResponseStream.WriteStructure(persistentGathering)
	}

	rmcResponseBody := rmcResponseStream.Bytes()
	fmt.Println(hex.EncodeToString(rmcResponseBody))

	rmcResponse := nex.NewRMCResponse(nexproto.MatchmakeExtensionProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchmakeExtensionMethodFindCommunityByGatheringId, rmcResponseBody)

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
