package nex

import (
	"github.com/PretendoNetwork/mario-kart-7/globals"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-common-go/matchmake-extension"
	matchmaking "github.com/PretendoNetwork/nex-protocols-common-go/matchmaking"
	match_making_ext "github.com/PretendoNetwork/nex-protocols-common-go/matchmaking-ext"
	nattraversal "github.com/PretendoNetwork/nex-protocols-common-go/nat-traversal"
	secure "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"

	nex_matchmake_extension_common "github.com/PretendoNetwork/mario-kart-7/nex/matchmake-extension/common"
)

func registerCommonSecureServerProtocols() {
	secure.NewCommonSecureConnectionProtocol(globals.SecureServer)

	nattraversal.NewCommonNATTraversalProtocol(globals.SecureServer)

	matchmaking.NewCommonMatchMakingProtocol(globals.SecureServer)

	match_making_ext.NewCommonMatchMakingExtProtocol(globals.SecureServer)

	matchmakeExtensionProtocol := matchmake_extension.NewCommonMatchmakeExtensionProtocol(globals.SecureServer)

	matchmakeExtensionProtocol.CleanupSearchMatchmakeSession(nex_matchmake_extension_common.CleanupSearchMatchmakeSession)
}
