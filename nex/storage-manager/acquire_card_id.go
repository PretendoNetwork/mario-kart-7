package nex_storage_manager

import (
	"math/rand"

	"github.com/PretendoNetwork/mario-kart-7/globals"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/v2/storage-manager"
)

func AcquireCardID(err error, packet nex.PacketInterface, callID uint32) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
	}

	cardID := types.NewPrimitiveU64(rand.Uint64())

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	cardID.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = storage_manager.ProtocolID
	rmcResponse.MethodID = storage_manager.MethodAcquireCardID
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
