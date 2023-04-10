package nex_matchmake_extension

import (
	"fmt"
	"math"
	"strconv"

	"github.com/PretendoNetwork/mario-kart-7-secure/database"
	"github.com/PretendoNetwork/mario-kart-7-secure/globals"
	nex "github.com/PretendoNetwork/nex-go"
	match_making "github.com/PretendoNetwork/nex-protocols-go/match-making"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
	"github.com/PretendoNetwork/nex-protocols-go/notifications"
)

func AutoMatchmake_Postpone(err error, client *nex.Client, callID uint32, matchmakeSession *match_making.MatchmakeSession, message string) {

	gameConfig := matchmakeSession.Attributes[2]
	fmt.Println(strconv.FormatUint(uint64(gameConfig), 2))
	fmt.Println("===== MATCHMAKE SESSION SEARCH =====")
	fmt.Println("COMMUNITY GID: " + strconv.Itoa((int)(matchmakeSession.Attributes[0])))
	fmt.Println("UPDATE: " + strconv.Itoa((int)(matchmakeSession.Attributes[2])))
	fmt.Println("REGION: " + strconv.Itoa((int)(matchmakeSession.Attributes[3])))
	//fmt.Println("REGION: " + regionList[matchmakeSession.Attributes[3]])
	fmt.Println("GAME MODE: " + strconv.Itoa((int)(matchmakeSession.GameMode)))
	//fmt.Println("GAME MODE: " + gameModes[matchmakeSession.GameMode])
	//fmt.Println("CC: " + ccList[gameConfig%0b111])
	gameConfig = gameConfig >> 12
	//fmt.Println("DLC MODE: " + dlcModes[matchmakeSession.Attributes[5]&0xF])
	//fmt.Println("ITEM MODE: " + itemModes[gameConfig%0b1111])
	gameConfig = gameConfig >> 8
	//fmt.Println("VEHICLE MODE: " + vehicleModes[gameConfig%0b11])
	gameConfig = gameConfig >> 4
	//fmt.Println("CONTROLLER MODE: " + controllerModes[gameConfig%0b11])
	//fmt.Println("HAVE GUEST PLAYER: " + strconv.FormatBool(false))

	fmt.Println("0: " + strconv.Itoa((int)(matchmakeSession.Attributes[0])))
	fmt.Println("2: " + strconv.Itoa((int)(matchmakeSession.Attributes[2])))
	fmt.Println("3: " + strconv.Itoa((int)(matchmakeSession.Attributes[3])))
	fmt.Println("5: " + strconv.Itoa((int)(matchmakeSession.Attributes[5])))

	gid := database.FindRoom(matchmakeSession.GameMode, true, matchmakeSession.Attributes[3], matchmakeSession.Attributes[2], uint32(1), matchmakeSession.Attributes[5]&0xF, matchmakeSession.Attributes[0])
	if gid == math.MaxUint32 {
		gid = database.NewRoom(client.PID(), matchmakeSession.GameMode, true, matchmakeSession.Attributes[3], matchmakeSession.Attributes[2], uint32(1), matchmakeSession.Attributes[5]&0xF, matchmakeSession.Attributes[0])
	}

	fmt.Println("GATHERING ID: " + strconv.Itoa((int)(gid)))

	database.AddPlayerToRoom(gid, client.PID(), uint32(1))

	hostpid, gamemode, _, _, update := database.GetRoomInfo(gid)

	matchmakeSession.Gathering.ID = gid
	matchmakeSession.Gathering.OwnerPID = hostpid
	matchmakeSession.Gathering.HostPID = hostpid
	matchmakeSession.GameMode = gamemode
	matchmakeSession.Attributes[2] = update

	rmcResponseStream := nex.NewStreamOut(globals.NEXServer)
	rmcResponseStream.WriteString("MatchmakeSession")
	lengthStream := nex.NewStreamOut(globals.NEXServer)
	lengthStream.WriteStructure(matchmakeSession.Gathering)
	lengthStream.WriteStructure(matchmakeSession)
	matchmakeSessionLength := uint32(len(lengthStream.Bytes()))
	rmcResponseStream.WriteUInt32LE(matchmakeSessionLength + 4)
	rmcResponseStream.WriteUInt32LE(matchmakeSessionLength)
	rmcResponseStream.WriteStructure(matchmakeSession.Gathering)
	rmcResponseStream.WriteStructure(matchmakeSession)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(matchmake_extension.ProtocolID, callID)
	rmcResponse.SetSuccess(matchmake_extension.MethodAutoMatchmake_Postpone, rmcResponseBody)

	var globalStation nex.StationURL
	fmt.Println(globals.NEXServer.FindClientFromPID(hostpid).ConnectionID())
	stationUrlStringsOrig := database.GetPlayerURLs(globals.NEXServer.FindClientFromPID(hostpid).ConnectionID())
	globalStation.FromString(stationUrlStringsOrig[1])
	address := client.Address().IP.String()
	port := strconv.Itoa(client.Address().Port)

	globalStation.SetAddress(address)
	globalStation.SetPort(port)

	globalStationURL := globalStation.EncodeToString()
	database.UpdatePlayerSessionURL(client.ConnectionID(), stationUrlStringsOrig[1], globalStationURL)

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

	rmcMessage := nex.NewRMCRequest()
	rmcMessage.SetProtocolID(notifications.ProtocolID)
	rmcMessage.SetCallID(0xffff0000 + callID)
	rmcMessage.SetMethodID(notifications.MethodProcessNotificationEvent)

	oEvent := notifications.NewNotificationEvent()
	oEvent.PIDSource = hostpid
	oEvent.Type = 3001 // New participant
	oEvent.Param1 = gid
	oEvent.Param2 = client.PID()

	stream := nex.NewStreamOut(globals.NEXServer)
	oEventBytes := oEvent.Bytes(stream)
	rmcMessage.SetParameters(oEventBytes)
	rmcMessageBytes := rmcMessage.Bytes()

	targetClient := globals.NEXServer.FindClientFromPID(uint32(hostpid))

	messagePacket, _ := nex.NewPacketV0(targetClient, nil)
	messagePacket.SetVersion(1)
	messagePacket.SetSource(0xA1)
	messagePacket.SetDestination(0xAF)
	messagePacket.SetType(nex.DataPacket)
	messagePacket.SetPayload(rmcMessageBytes)

	messagePacket.AddFlag(nex.FlagNeedsAck)
	messagePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(messagePacket)

	messagePacket, _ = nex.NewPacketV0(client, nil)
	messagePacket.SetVersion(1)
	messagePacket.SetSource(0xA1)
	messagePacket.SetDestination(0xAF)
	messagePacket.SetType(nex.DataPacket)
	messagePacket.SetPayload(rmcMessageBytes)

	messagePacket.AddFlag(nex.FlagNeedsAck)
	messagePacket.AddFlag(nex.FlagReliable)

	globals.NEXServer.Send(messagePacket)
}
