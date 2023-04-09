package nex

import (
	"github.com/PretendoNetwork/mario-kart-7-secure/globals"
	match_making_ext "github.com/PretendoNetwork/nex-protocols-go/match-making-ext"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/matchmake-extension"
	"github.com/PretendoNetwork/nex-protocols-go/ranking"
	"github.com/PretendoNetwork/nex-protocols-go/utility"

	nex_match_making_ext "github.com/PretendoNetwork/mario-kart-7-secure/nex/match-making-ext"
	nex_matchmake_extension "github.com/PretendoNetwork/mario-kart-7-secure/nex/matchmake-extension"
	nex_ranking "github.com/PretendoNetwork/mario-kart-7-secure/nex/ranking"
	nex_utility "github.com/PretendoNetwork/mario-kart-7-secure/nex/utility"
)

func registerNEXProtocols() {
	matchmakeExtensionProtocol := matchmake_extension.NewMatchmakeExtensionProtocol(globals.NEXServer)

	matchmakeExtensionProtocol.CloseParticipation(nex_matchmake_extension.CloseParticipation)
	matchmakeExtensionProtocol.AutoMatchmakeWithParam_Postpone(nex_matchmake_extension.AutoMatchmakeWithParam_Postpone)
	matchmakeExtensionProtocol.AutoMatchmake_Postpone(nex_matchmake_extension.AutoMatchmake_Postpone)
	matchmakeExtensionProtocol.GetPlayingSession(nex_matchmake_extension.GetPlayingSession)
	matchmakeExtensionProtocol.CreateCommunity(nex_matchmake_extension.CreateCommunity)
	matchmakeExtensionProtocol.FindCommunityByGatheringID(nex_matchmake_extension.FindCommunityByGatheringID)
	matchmakeExtensionProtocol.FindOfficialCommunity(nex_matchmake_extension.FindOfficialCommunity)
	matchmakeExtensionProtocol.FindCommunityByParticipant(nex_matchmake_extension.FindCommunityByParticipant)
	matchmakeExtensionProtocol.GetSimpleCommunity(nex_matchmake_extension.GetSimpleCommunity)
	matchmakeExtensionProtocol.UpdateProgressScore(nex_matchmake_extension.UpdateProgressScore)
	matchmakeExtensionProtocol.CreateMatchmakeSessionWithParam(nex_matchmake_extension.CreateMatchmakeSessionWithParam)
	matchmakeExtensionProtocol.JoinMatchmakeSessionWithParam(nex_matchmake_extension.JoinMatchmakeSessionWithParam)

	matchMakingExtProtocol := match_making_ext.NewMatchMakingExtProtocol(globals.NEXServer)

	matchMakingExtProtocol.EndParticipation(nex_match_making_ext.EndParticipation)

	rankingProtocol := ranking.NewRankingProtocol(globals.NEXServer)

	rankingProtocol.UploadCommonData(nex_ranking.UploadCommonData)

	utilityProtocol := utility.NewUtilityProtocol(globals.NEXServer)

	utilityProtocol.GetAssociatedNexUniqueIDWithMyPrincipalID(nex_utility.GetAssociatedNexUniqueIDWithMyPrincipalID)
	utilityProtocol.AssociateNexUniqueIDsWithMyPrincipalID(nex_utility.AssociateNexUniqueIDsWithMyPrincipalID)
}
