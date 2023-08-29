package nex

import (
	"github.com/PretendoNetwork/mario-kart-7/globals"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-common-go/matchmake-extension"
	matchmaking "github.com/PretendoNetwork/nex-protocols-common-go/matchmaking"
	matchmaking_ext "github.com/PretendoNetwork/nex-protocols-common-go/matchmaking-ext"
	nattraversal "github.com/PretendoNetwork/nex-protocols-common-go/nat-traversal"
	secure "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"

	nex_matchmake_extension "github.com/PretendoNetwork/mario-kart-7/nex/matchmake-extension"
)

func registerCommonSecureServerProtocols() {
	secureConnectionCommonProtocol := secure.NewCommonSecureConnectionProtocol(globals.SecureServer)
	_ = secureConnectionCommonProtocol

	natTraversalCommonProtocol := nattraversal.NewCommonNATTraversalProtocol(globals.SecureServer)
	_ = natTraversalCommonProtocol

	matchMakingCommonProtocol := matchmaking.NewCommonMatchMakingProtocol(globals.SecureServer)
	_ = matchMakingCommonProtocol

	matchMakingCommonExtProtocol := matchmaking_ext.NewCommonMatchMakingExtProtocol(globals.SecureServer)
	_ = matchMakingCommonExtProtocol

	matchmakeExtensionProtocol := matchmake_extension.NewCommonMatchmakeExtensionProtocol(globals.SecureServer, "")

	matchmakeExtensionProtocol.CleanupSearchMatchmakeSession(nex_matchmake_extension.CleanupSearchMatchmakeSession)
}
