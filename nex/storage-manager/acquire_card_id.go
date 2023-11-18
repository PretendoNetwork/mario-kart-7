package nex_storage_manager

import (
	"math/rand"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/storage-manager"
)

func AcquireCardID(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.Core.Unknown
	}

	cardID := rand.Uint64()

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteUInt64LE(cardID)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(rmcResponseBody)
	rmcResponse.ProtocolID = storage_manager.ProtocolID
	rmcResponse.MethodID = storage_manager.MethodAcquireCardID
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
