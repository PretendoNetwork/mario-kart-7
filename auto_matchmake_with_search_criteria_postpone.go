package main

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var regionList = []string{"Worldwide", "Japan", "United States", "Europe", "Korea", "China", "Taiwan"}
var gameModes = []string{"VS Race", "Battle"}
var ccList = []string{"Unk", "200cc", "50cc", "100cc", "150cc", "Mirror", "BattleCC"}
var itemModes = []string{"Unk1", "Unk2", "Unk3", "Unk4", "Unk5", "Normal", "Unk7", "All Items", "Shells Only", "Bananas Only", "Mushrooms Only", "Bob-ombs Only", "No Items", "No Items or Coins", "Frantic"}
var vehicleModes = []string{"All Vehicles", "Karts Only", "Bikes Only"}
var controllerModes = []string{"Unk", "Tilt Only", "All Controls"}
var dlcModes = []string{"No DLC", "DLC Pack 1 Only", "DLC Pack 2 Only", "Both DLC Packs"}

func autoMatchmakeWithSearchCriteria_Postpone(err error, client *nex.Client, callID uint32, searchCriteria []*nexproto.MatchmakeSessionSearchCriteria, matchmakeSession *nexproto.MatchmakeSession, message string) {

	gameConfig := matchmakeSession.Attributes[2]
	fmt.Println(strconv.FormatUint(uint64(gameConfig), 2))
	fmt.Println("===== MATCHMAKE SESSION JOIN =====")
	fmt.Println("REGION: " + regionList[matchmakeSession.Attributes[3]])
	fmt.Println("GAME MODE: " + gameModes[matchmakeSession.GameMode])
	fmt.Println("CC: " + ccList[gameConfig%0b111])
	gameConfig = gameConfig >> 12
	fmt.Println("DLC MODE: " + dlcModes[matchmakeSession.Attributes[5]&0xF])
	fmt.Println("ITEM MODE: " + itemModes[gameConfig%0b1111])
	gameConfig = gameConfig >> 8
	fmt.Println("VEHICLE MODE: " + vehicleModes[gameConfig%0b11])
	gameConfig = gameConfig >> 4
	fmt.Println("CONTROLLER MODE: " + controllerModes[gameConfig%0b11])
	fmt.Println("HAVE GUEST PLAYER: " + strconv.FormatBool(searchCriteria[0].VacantParticipants > 1))

	gid := findRoom(matchmakeSession.GameMode, true, matchmakeSession.Attributes[3], matchmakeSession.Attributes[2], uint32(searchCriteria[0].VacantParticipants), matchmakeSession.Attributes[5]&0xF)
	if gid == math.MaxUint32 {
		gid = newRoom(client.PID(), matchmakeSession.GameMode, true, matchmakeSession.Attributes[3], matchmakeSession.Attributes[2], uint32(searchCriteria[0].VacantParticipants), matchmakeSession.Attributes[5]&0xF)
	}

	addPlayerToRoom(gid, client.PID(), uint32(searchCriteria[0].VacantParticipants))

	hostpid, gamemode, region, gconfig, dlcmode := getRoomInfo(gid)
	sessionKey := "00000000000000000000000000000000"

	matchmakeSession.Gathering.ID = gid
	matchmakeSession.Gathering.OwnerPID = hostpid
	matchmakeSession.Gathering.HostPID = hostpid
	matchmakeSession.Gathering.MinimumParticipants = 1
	matchmakeSession.SessionKey = []byte(sessionKey)
	matchmakeSession.GameMode = gamemode
	matchmakeSession.Attributes[3] = region
	matchmakeSession.Attributes[2] = gconfig
	matchmakeSession.Attributes[5] = dlcmode

	rmcResponseStream := nex.NewStreamOut(nexServer)
	rmcResponseStream.WriteString("MatchmakeSession")
	lengthStream := nex.NewStreamOut(nexServer)
	lengthStream.WriteStructure(matchmakeSession.Gathering)
	lengthStream.WriteStructure(matchmakeSession)
	matchmakeSessionLength := uint32(len(lengthStream.Bytes()))
	rmcResponseStream.WriteUInt32LE(matchmakeSessionLength + 4)
	rmcResponseStream.WriteUInt32LE(matchmakeSessionLength)
	rmcResponseStream.WriteStructure(matchmakeSession.Gathering)
	rmcResponseStream.WriteStructure(matchmakeSession)

	rmcResponseBody := rmcResponseStream.Bytes()

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

func endParticipation(err error, client *nex.Client, callID uint32, idGathering uint32, strMessage string) {
	removePlayerFromRoom(idGathering, client.PID())

	returnval := []byte{0x1}

	rmcResponse := nex.NewRMCResponse(nexproto.MatchMakingExtProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchMakingExtMethodEndParticipation, returnval)

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

func sendReport(err error, client *nex.Client, callID uint32, reportID uint32, reportData []byte) {
	fmt.Println("Report ID: " + strconv.Itoa(int(reportID)))
	fmt.Println("Report Data: " + hex.EncodeToString(reportData))

	rmcResponse := nex.NewRMCResponse(nexproto.SecureProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.SecureMethodSendReport, nil)

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

func updateSessionHostV1(err error, client *nex.Client, callID uint32, gid uint32) {

	updateRoomHost(gid, client.PID())

	rmcResponse := nex.NewRMCResponse(nexproto.MatchMakingProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.MatchMakingMethodUpdateSessionHostV1, nil)

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
