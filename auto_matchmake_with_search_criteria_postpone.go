package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
	"fmt"
	"encoding/hex"
)

func autoMatchmakeWithSearchCriteria_Postpone(err error, client *nex.Client, callID uint32, matchmakeSession *nexproto.MatchmakeSession, message string) {
	var foundSession *MatchmakingData
	var foundIndex int
	var element *MatchmakingData
	rmcResponseStream := nex.NewStreamOut(nexServer)
	for foundIndex, element = range MatchmakingState[1:] {
		if(uint16(len(element.clients)) < element.matchmakeSession.Gathering.MaximumParticipants ){
			foundSession = element
		}
	}
	sessionKey := "00000000000000000000000000000000"
	if(foundSession == nil){
		foundSession = new(MatchmakingData)
		foundIndex = len(MatchmakingState)
		matchmakeSession.Gathering.ID = uint32(foundIndex)
		matchmakeSession.Gathering.OwnerPID = client.PID()
		matchmakeSession.Gathering.HostPID = client.PID()
		matchmakeSession.Gathering.MinimumParticipants = 1
		foundSession.matchmakeSession = matchmakeSession
		foundSession.clients = make([]*nex.Client, 0, 0)
		matchmakeSession.SessionKey = []byte(sessionKey)
	}
	fmt.Println(foundSession.matchmakeSession.Gathering.OwnerPID)
	foundSession.clients = append(foundSession.clients, client)
	MatchmakingState = append(MatchmakingState, foundSession)
	rmcResponseStream.WriteString("MatchmakeSession")
	lengthStream := nex.NewStreamOut(nexServer)
	lengthStream.WriteStructure(foundSession.matchmakeSession.Gathering)
	lengthStream.WriteStructure(foundSession.matchmakeSession)
	matchmakeSessionLength := uint32(len(lengthStream.Bytes()))
	rmcResponseStream.WriteUInt32LE(matchmakeSessionLength+4)
	rmcResponseStream.WriteUInt32LE(matchmakeSessionLength)
	rmcResponseStream.WriteStructure(foundSession.matchmakeSession.Gathering)
	rmcResponseStream.WriteStructure(foundSession.matchmakeSession)

	rmcResponseBody := rmcResponseStream.Bytes()
	fmt.Println(hex.EncodeToString(rmcResponseBody))

	rmcResponse := nex.NewRMCResponse(nexproto.MatchmakeExtensionProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchmakeExtensionMethodAutoMatchmakeWithSearchCriteria_Postpone, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
