package nex_storage_manager

import (
	"database/sql"

	"github.com/PretendoNetwork/mario-kart-7/database"
	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/storage-manager"
)

func ActivateWithCardID(err error, packet nex.PacketInterface, callID uint32, unknown uint8, cardID uint64) (*nex.RMCMessage, uint32) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.Errors.Core.InvalidArgument
	}

	client := packet.Sender()

	// * It's not guaranteed that the client will call AcquireCardID,
	// * because that method is only called the first time the client
	// * goes online, and it may have used official servers previously.
	// *
	// * To workaround this, we ignore the card ID stuff and get the
	// * unique ID using the PID
	var firstTime bool
	uniqueID, err := database.GetUniqueIDByOwnerPID(client.PID().LegacyValue())
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nil, nex.Errors.Core.Unknown
	}

	if err == sql.ErrNoRows {
		uniqueID, err = database.InsertCommonDataByOwnerPID(client.PID().LegacyValue())
		if err != nil {
			globals.Logger.Critical(err.Error())
			return nil, nex.Errors.Core.Unknown
		}

		firstTime = true
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteUInt32LE(uniqueID)
	rmcResponseStream.WriteBool(firstTime)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(globals.SecureServer, rmcResponseBody)
	rmcResponse.ProtocolID = storage_manager.ProtocolID
	rmcResponse.MethodID = storage_manager.MethodActivateWithCardID
	rmcResponse.CallID = callID

	return rmcResponse, 0
}
