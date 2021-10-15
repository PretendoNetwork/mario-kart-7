package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

type MatchmakingData struct {
	matchmakeSession    *nexproto.MatchmakeSession
	clients             []*nex.Client
}

var nexServer *nex.Server
var secureServer *nexproto.SecureProtocol
var MatchmakingState []*MatchmakingData

func main() {
	MatchmakingState = append(MatchmakingState, nil)

	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(1)
	nexServer.SetNexVersion(30500)
	nexServer.SetKerberosKeySize(32)
	nexServer.SetAccessKey("25dbf96a")

	nexServer.On("Data", func(packet *nex.PacketV1) {
		request := packet.RMCRequest()

		fmt.Println("==MK8 - Secure==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("=================")
	})

	secureServer = nexproto.NewSecureProtocol(nexServer)
	natTraversalProtocolServer := nexproto.NewNatTraversalProtocol(nexServer)
	utilityProtocolServer := nexproto.NewUtilityProtocol(nexServer)
	matchmakeExtensionProtocolServer := nexproto.NewMatchmakeExtensionProtocol(nexServer)
	matchMakingProtocolServer := nexproto.NewMatchMakingProtocol(nexServer)
	rankingProtocolServer := nexproto.NewRankingProtocol(nexServer)

	//needed for the datastore method MK7 contacts when first going online (just needs a response of some kind)
	dataStorePrococolServer := nexproto.NewDataStoreProtocol(nexServer)
	_ = dataStorePrococolServer

	// Handle PRUDP CONNECT packet (not an RMC method)
	nexServer.On("Connect", connect)

	secureServer.Register(register)
	secureServer.ReplaceURL(replaceURL)

	natTraversalProtocolServer.RequestProbeInitiationExt(requestProbeInitiationExt)
	natTraversalProtocolServer.ReportNatProperties(reportNatProperties)

	utilityProtocolServer.GetAssociatedNexUniqueIdWithMyPrincipalId(getAssociatedNexUniqueIdWithMyPrincipalId)

	matchmakeExtensionProtocolServer.AutoMatchmakeWithSearchCriteria_Postpone(autoMatchmakeWithSearchCriteria_Postpone)

	matchMakingProtocolServer.GetSessionURLs(getSessionURLs)
	
	rankingProtocolServer.UploadCommonData(uploadCommonData)

	nexServer.Listen(":60003")
}
