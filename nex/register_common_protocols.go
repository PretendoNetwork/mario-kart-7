package nex

import (
	"github.com/PretendoNetwork/mario-kart-7-secure/database"
	"github.com/PretendoNetwork/mario-kart-7-secure/globals"
	nexmatchmaking "github.com/PretendoNetwork/nex-protocols-common-go/matchmaking"
	nexnattraversal "github.com/PretendoNetwork/nex-protocols-common-go/nat-traversal"
	nexsecure "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
)

func registerCommonProtocols() {
	secureConnectionCommonProtocol := nexsecure.NewCommonSecureConnectionProtocol(globals.NEXServer)
	secureConnectionCommonProtocol.AddConnection(database.AddPlayerSession)
	secureConnectionCommonProtocol.UpdateConnection(database.UpdatePlayerSessionAll)
	secureConnectionCommonProtocol.DoesConnectionExist(database.DoesSessionExist)
	secureConnectionCommonProtocol.ReplaceConnectionUrl(database.UpdatePlayerSessionURL)

	natTraversalCommonProtocol := nexnattraversal.InitNatTraversalProtocol(globals.NEXServer)
	nexnattraversal.GetConnectionUrls(database.GetPlayerURLs)
	nexnattraversal.ReplaceConnectionUrl(database.UpdatePlayerSessionURL)
	_ = natTraversalCommonProtocol

	matchMakingCommonProtocol := nexmatchmaking.InitMatchmakingProtocol(globals.NEXServer)
	nexmatchmaking.GetConnectionUrls(database.GetPlayerURLs)
	nexmatchmaking.UpdateRoomHost(database.UpdateRoomHost)
	nexmatchmaking.DestroyRoom(database.DestroyRoom)
	nexmatchmaking.GetRoomInfo(database.GetRoomInfo)
	nexmatchmaking.GetRoomPlayers(database.GetRoomPlayers)
	_ = matchMakingCommonProtocol
}
