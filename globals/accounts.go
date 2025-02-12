package globals

import (
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

var AuthenticationServerAccount *nex.Account
var SecureServerAccount *nex.Account

func AccountDetailsByPID(pid types.PID) (*nex.Account, *nex.Error) {
	if pid.Equals(AuthenticationServerAccount.PID) {
		return AuthenticationServerAccount, nil
	}

	if pid.Equals(SecureServerAccount.PID) {
		return SecureServerAccount, nil
	}

	password, errorCode := PasswordFromPID(pid)
	if errorCode != 0 {
		return nil, nex.NewError(errorCode, "Failed to get password from PID")
	}

	account := nex.NewAccount(pid, strconv.Itoa(int(pid)), password)

	return account, nil
}

func AccountDetailsByUsername(username string) (*nex.Account, *nex.Error) {
	if username == AuthenticationServerAccount.Username {
		return AuthenticationServerAccount, nil
	}

	if username == SecureServerAccount.Username {
		return SecureServerAccount, nil
	}

	pidInt, err := strconv.Atoi(username)
	if err != nil {
		return nil, nex.NewError(nex.ResultCodes.RendezVous.InvalidUsername, "Invalid username")
	}

	pid := types.NewPID(uint64(pidInt))

	password, errorCode := PasswordFromPID(pid)
	if errorCode != 0 {
		return nil, nex.NewError(errorCode, "Failed to get password from PID")
	}

	account := nex.NewAccount(pid, username, password)

	return account, nil
}
