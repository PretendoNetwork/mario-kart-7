package main

import (
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func reportNatProperties(err error, client *nex.Client, callID uint32, natm uint32, natf uint32, rtt uint32) {
	stationUrlsStrings := getPlayerUrls(client.PID())
	stationUrls := make([]nex.StationURL, len(stationUrlsStrings))
	pid := strconv.FormatUint(uint64(client.PID()), 10)
	rvcid := strconv.FormatUint(uint64(client.ConnectionId()), 10)

	for i := 0; i < len(stationUrlsStrings); i++ {
		stationUrls[i] = *nex.NewStationURL(stationUrlsStrings[i])
		if stationUrls[i].Type() == "3" {
			natm_s := strconv.FormatUint(uint64(natm), 10)
			natf_s := strconv.FormatUint(uint64(natf), 10)
			stationUrls[i].SetNatm(&natm_s)
			stationUrls[i].SetNatf(&natf_s)
		}
		stationUrls[i].SetPid(&pid)
		stationUrls[i].SetRVCID(&rvcid)
		updatePlayerSessionUrl(client.PID(), stationUrlsStrings[i], stationUrls[i].EncodeToString())
	}

	rmcResponse := nex.NewRMCResponse(nexproto.NatTraversalProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.NatTraversalMethodReportNatProperties, nil)

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
