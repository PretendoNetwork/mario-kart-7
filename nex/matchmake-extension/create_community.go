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

func CreateCommunity(err error, client *nex.Client, callID uint32, community *match_making.PersistentGathering, strMessage string) {
	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)
	gid := database.NewCommunity(client.PID(), community.M_CommunityType, community.M_Password, community.M_Attribs, community.M_ApplicationBuffer, community.M_ParticipationStartDate.Value(), community.M_ParticipationEndDate.Value())

	rmcResponseStream.WriteUInt32LE(gid)
	rmcResponseBody := rmcResponseStream.Bytes()
	fmt.Println(hex.EncodeToString(rmcResponseBody))

	rmcResponse := nex.NewRMCResponse(matchmake_extension.ProtocolID, callID)
	rmcResponse.SetSuccess(matchmake_extension.MethodCreateCommunity, rmcResponseBody)

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
