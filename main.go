package main

import (
	"fmt"
	"time"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

type MatchmakingData struct {
	matchmakeSession *nexproto.MatchmakeSession
	clients          []*nex.Client
}

var nexServer *nex.Server
var secureServer *nexproto.SecureProtocol
var MatchmakingState []*MatchmakingData
var config *ServerConfig

func main() {
	MatchmakingState = append(MatchmakingState, nil)

	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(1)
	nexServer.SetNexVersion(config.NexVersion)
	nexServer.SetKerberosKeySize(32)
	nexServer.SetAccessKey(config.AccessKey)

	nexServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==MK8 - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("=================")
	})

	secureServer = nexproto.NewSecureProtocol(nexServer)
	natTraversalProtocolServer := nexproto.NewNatTraversalProtocol(nexServer)
	matchmakeExtensionProtocolServer := nexproto.NewMatchmakeExtensionProtocol(nexServer)
	matchMakingProtocolServer := nexproto.NewMatchMakingProtocol(nexServer)
	matchMakingExtProtocolServer := nexproto.NewMatchMakingExtProtocol(nexServer)
	rankingProtocolServer := nexproto.NewRankingProtocol(nexServer)

	// have datastore available if called, but respond as unimplemented
	dataStorePrococolServer := nexproto.NewDataStoreProtocol(nexServer)
	_ = dataStorePrococolServer

	// Handle PRUDP CONNECT packet (not an RMC method)
	nexServer.On("Connect", connect)

	secureServer.Register(register)
	secureServer.ReplaceURL(replaceURL)
	secureServer.SendReport(sendReport)

	natTraversalProtocolServer.RequestProbeInitiationExt(requestProbeInitiationExt)
	natTraversalProtocolServer.ReportNatProperties(reportNatProperties)

	matchmakeExtensionProtocolServer.AutoMatchmakeWithSearchCriteria_Postpone(autoMatchmakeWithSearchCriteria_Postpone)

	matchMakingProtocolServer.GetSessionURLs(getSessionURLs)
	matchMakingProtocolServer.UpdateSessionHostV1(updateSessionHostV1)

	matchMakingExtProtocolServer.EndParticipation(endParticipation)

	rankingProtocolServer.UploadCommonData(uploadCommonData)

	nexServer.Listen(":" + config.ServerPort)
}

// Modified version of https://gist.github.com/ryanfitz/4191392

// will eventually be used to occasionally check for disconnected clients
// so as to clean their session info out of the database
func doEvery(d time.Duration, f func()) {
	for x := range time.Tick(d) {
		x = x
		f()
	}
}
