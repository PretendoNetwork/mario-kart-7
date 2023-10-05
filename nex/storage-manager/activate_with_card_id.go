package nex_storage_manager

import (
	"database/sql"

	"github.com/PretendoNetwork/mario-kart-7/database"
	"github.com/PretendoNetwork/mario-kart-7/globals"
	nex "github.com/PretendoNetwork/nex-go"
	storage_manager "github.com/PretendoNetwork/nex-protocols-go/storage-manager"
)

func ActivateWithCardID(err error, client *nex.Client, callID uint32, unknown uint8, cardID uint64) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.Errors.Core.InvalidArgument
	}

	// * It's not guaranteed that the client will call AcquireCardID,
	// * because that method is only called the first time the client
	// * goes online, and it may have used official servers previously.
	// *
	// * To workaround this, we ignore the card ID stuff and get the
	// * unique ID using the PID
	var firstTime bool
	uniqueID, err := database.GetUniqueIDByOwnerPID(client.PID())
	if err != nil && err != sql.ErrNoRows {
		globals.Logger.Critical(err.Error())
		return nex.Errors.Core.Unknown
	}

	if err == sql.ErrNoRows {
		uniqueID, err = database.InsertCommonDataByOwnerPID(client.PID())
		if err != nil {
			globals.Logger.Critical(err.Error())
			return nex.Errors.Core.Unknown
		}

		firstTime = true
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteUInt32LE(uniqueID)
	rmcResponseStream.WriteBool(firstTime)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(storage_manager.ProtocolID, callID)
	rmcResponse.SetSuccess(storage_manager.MethodActivateWithCardID, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0
}
