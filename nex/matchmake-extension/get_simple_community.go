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

func GetSimpleCommunity(err error, client *nex.Client, callID uint32, gatheringIdList []uint32) {
	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)
	lstSimpleCommunityList := make([]*match_making.SimpleCommunity, 0)
	for _, gid := range gatheringIdList {
		fmt.Println("gid1", gid)
		if database.CommunityExists(gid) {
			fmt.Println("gid2", gid)
			_, _, _, _, _, _, _, sessions, _ := database.GetCommunityInfo(gid)
			simpleCommunity := match_making.NewSimpleCommunity()
			simpleCommunity.M_GatheringID = gid
			simpleCommunity.M_MatchmakeSessionCount = sessions
			lstSimpleCommunityList = append(lstSimpleCommunityList, simpleCommunity)
		}
	}

	rmcResponseStream.WriteUInt32LE(uint32(len(lstSimpleCommunityList)))
	for _, simpleCommunity := range lstSimpleCommunityList {
		rmcResponseStream.WriteStructure(simpleCommunity)
	}
	rmcResponseBody := rmcResponseStream.Bytes()
	fmt.Println(hex.EncodeToString(rmcResponseBody))

	rmcResponse := nex.NewRMCResponse(matchmake_extension.ProtocolID, callID)
	rmcResponse.SetSuccess(matchmake_extension.MethodGetSimpleCommunity, rmcResponseBody)

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
