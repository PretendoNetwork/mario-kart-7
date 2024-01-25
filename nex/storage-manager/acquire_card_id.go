package nex_storage_manager

import (
	"math/rand"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-go/types"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/storage-manager"
)

func AcquireCardID(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.ResultCodes.Core.Unknown
	}

	cardID := types.NewPrimitiveU64(rand.Uint64())

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer)

	cardID.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, rmcResponseBody)
	rmcResponse.ProtocolID = storage_manager.ProtocolID
	rmcResponse.MethodID = storage_manager.MethodAcquireCardID
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
