package nex_matchmake_extension_common

import (
	matchmaking_types "github.com/PretendoNetwork/nex-protocols-go/match-making/types"
)

func CleanupSearchMatchmakeSession(matchmakeSession *matchmaking_types.MatchmakeSession) {
	// Cleanup VR
	matchmakeSession.Attributes[1] = 0
}
