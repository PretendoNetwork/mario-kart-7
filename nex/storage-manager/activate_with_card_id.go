package nex_storage_manager

import (
	"database/sql"

	"github.com/PretendoNetwork/mario-kart-7/database"
	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/v2/storage-manager"
)

func ActivateWithCardID(err error, packet nex.PacketInterface, callID uint32, unknown types.UInt8, cardID types.UInt64) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, err.Error())
	}

	client := packet.Sender()

	uniqueID := types.NewUInt32(0)
	firstTime := types.NewBool(false)

	// * It's not guaranteed that the client will call AcquireCardID,
	// * because that method is only called the first time the client
	// * goes online, and it may have used official servers previously.
	// *
	// * To workaround this, we ignore the card ID stuff and get the
	// * unique ID using the PID
	rawUniqueID, err := database.GetUniqueIDByOwnerPID(uint32(client.PID()))
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
	}

	if err == sql.ErrNoRows {
		rawUniqueID, err = database.InsertCommonDataByOwnerPID(uint32(client.PID()))
		if err != nil {
			globals.Logger.Critical(err.Error())
			return nil, nex.NewError(nex.ResultCodes.Core.Unknown, err.Error())
		}

		firstTime = true
	}

	uniqueID = types.NewUInt32(rawUniqueID)

	rmcResponseStream := nex.NewByteStreamOut(globals.SecureServer.LibraryVersions, globals.SecureServer.ByteStreamSettings)

	uniqueID.WriteTo(rmcResponseStream)
	firstTime.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureEndpoint, rmcResponseBody)
	rmcResponse.ProtocolID = storage_manager.ProtocolID
	rmcResponse.MethodID = storage_manager.MethodActivateWithCardID
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
