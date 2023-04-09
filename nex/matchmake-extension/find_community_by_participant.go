package nex_matchmake_extension

import (
	"encoding/hex"
	"fmt"

	"github.com/PretendoNetwork/mario-kart-7-secure/database"
	"github.com/PretendoNetwork/mario-kart-7-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	match_making "github.com/PretendoNetwork/nex-protocols-go/match-making"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
)

func FindCommunityByParticipant(err error, client *nex.Client, callID uint32, pid uint32, resultRange *nex.ResultRange) {
	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)
	lstCommunity := make([]*match_making.PersistentGathering, 0)
	lstGid := database.FindCommunitiesWithParticipant(pid)
	for _, gid := range lstGid {
		fmt.Println("gid1", gid)
		if database.CommunityExists(gid) {
			fmt.Println("gid2", gid)
			host, communityType, password, attribs, application_buffer, start_date, end_date, sessions, participationCount := database.GetCommunityInfo(gid)
			persistentGathering := match_making.NewPersistentGathering()
			persistentGathering.Gathering.ID = gid
			persistentGathering.Gathering.OwnerPID = host
			persistentGathering.Gathering.HostPID = host
			persistentGathering.Gathering.MinimumParticipants = 1
			persistentGathering.Gathering.MaximumParticipants = 8
			persistentGathering.M_CommunityType = communityType
			persistentGathering.M_Password = password
			persistentGathering.M_Attribs = attribs
			persistentGathering.M_ApplicationBuffer = application_buffer
			persistentGathering.M_ParticipationStartDate = nex.NewDateTime(start_date)
			persistentGathering.M_ParticipationEndDate = nex.NewDateTime(end_date)
			persistentGathering.M_MatchmakeSessionCount = sessions
			persistentGathering.M_ParticipationCount = participationCount
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

	rmcResponse := nex.NewRMCResponse(matchmake_extension.ProtocolID, callID)
	rmcResponse.SetSuccess(matchmake_extension.MethodFindCommunityByParticipant, rmcResponseBody)

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
