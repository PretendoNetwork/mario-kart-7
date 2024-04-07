package nex

import (
	"github.com/PretendoNetwork/mario-kart-7/globals"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/v2/storage-manager"

	nex_storage_manager "github.com/PretendoNetwork/mario-kart-7/nex/storage-manager"
)

func registerSecureServerNEXProtocols() {
	datastoreProtocol := datastore.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(datastoreProtocol)

	// TODO - Add legacy ranking protocol!
	rankingProtocol := ranking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(rankingProtocol)

	storageManagerProtocol := storage_manager.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(storageManagerProtocol)

	storageManagerProtocol.AcquireCardID = nex_storage_manager.AcquireCardID
	storageManagerProtocol.ActivateWithCardID = nex_storage_manager.ActivateWithCardID
}
