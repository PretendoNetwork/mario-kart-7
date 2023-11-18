package nex

import (
	"github.com/PretendoNetwork/mario-kart-7/globals"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	ranking "github.com/PretendoNetwork/nex-protocols-go/ranking"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/storage-manager"

	nex_storage_manager "github.com/PretendoNetwork/mario-kart-7/nex/storage-manager"
)

func registerSecureServerNEXProtocols() {
	_ = datastore.NewProtocol(globals.SecureServer)

	// TODO - Add legacy ranking protocol!
	_ = ranking.NewProtocol(globals.SecureServer)

	storageManagerProtocol := storage_manager.NewProtocol(globals.SecureServer)

	storageManagerProtocol.AcquireCardID = nex_storage_manager.AcquireCardID
	storageManagerProtocol.ActivateWithCardID = nex_storage_manager.ActivateWithCardID
}
