package nex_storage_manager

import (
	"database/sql"

	"github.com/PretendoNetwork/mario-kart-7/database"
	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/nex-go/types"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/storage-manager"
)

func ActivateWithCardID(err error, packet nex.PacketInterface, callID uint32, unknown *types.PrimitiveU8, cardID *types.PrimitiveU64) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.ResultCodes.Core.InvalidArgument
	}

	client := packet.Sender()

	uniqueID := types.NewPrimitiveU32(0)
	firstTime := types.NewPrimitiveBool(false)

	// * It's not guaranteed that the client will call AcquireCardID,
	// * because that method is only called the first time the client
	// * goes online, and it may have used official servers previously.
	// *
	// * To workaround this, we ignore the card ID stuff and get the
	// * unique ID using the PID
	uniqueID.Value, err = database.GetUniqueIDByOwnerPID(client.PID().LegacyValue())
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nil, nex.ResultCodes.Core.Unknown
	}

	if err == sql.ErrNoRows {
		uniqueID.Value, err = database.InsertCommonDataByOwnerPID(client.PID().LegacyValue())
		if err != nil {
			globals.Logger.Critical(err.Error())
			return nil, nex.ResultCodes.Core.Unknown
		}

		firstTime.Value = true
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer)

	uniqueID.WriteTo(rmcResponseStream)
	firstTime.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, rmcResponseBody)
	rmcResponse.ProtocolID = storage_manager.ProtocolID
	rmcResponse.MethodID = storage_manager.MethodActivateWithCardID
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
