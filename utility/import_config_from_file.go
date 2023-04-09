package utility

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/PretendoNetwork/mario-kart-7-secure/globals"
	"github.com/PretendoNetwork/mario-kart-7-secure/types"
)

func ImportConfigFromFile(path string) (*types.ServerConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	indexes := strings.Split(string(data), "\n")
	config := &types.ServerConfig{
		ServerName:      "server",
		KerberosKeySize: 32,
	}
	for i := 0; i < len(indexes); i++ {
		index := strings.Split(indexes[i], "=")
		if len(index) != 2 {
			continue
		}
		switch index[0] {
		case "ServerName":
			globals.Config.ServerName = index[1]
			break
		case "ServerPort":
			globals.Config.ServerPort = index[1]
			break
		case "PrudpVersion":
			globals.Config.PrudpVersion, err = strconv.Atoi(index[1])
			if err != nil {
				return nil, err
			}
			break
		case "SignatureVersion":
			globals.Config.SignatureVersion, err = strconv.Atoi(index[1])
			if err != nil {
				return nil, err
			}
			break
		case "KerberosKeySize":
			globals.Config.KerberosKeySize, err = strconv.Atoi(index[1])
			if err != nil {
				return nil, err
			}
			break
		case "AccessKey":
			globals.Config.AccessKey = index[1]
			break
		case "NexVersion":
			globals.Config.NexVersion, err = strconv.Atoi(index[1])
			if err != nil {
				return nil, err
			}
			break
		case "AccountDatabaseIP":
			globals.Config.AccountDatabaseIP = index[1]
			break
		case "AccountDatabasePort":
			globals.Config.AccountDatabasePort = index[1]
			break
		case "NEXDatabaseIP":
			globals.Config.NEXDatabaseIP = index[1]
			break
		case "NEXDatabasePort":
			globals.Config.NEXDatabasePort = index[1]
			break
		case "DatabaseUseAuth":
			globals.Config.DatabaseUseAuth = (index[1] == "true")
			break
		case "DatabaseUsername":
			globals.Config.DatabaseUsername = index[1]
			break
		case "DatabasePassword":
			globals.Config.DatabasePassword = index[1]
			break
		case "AccountDatabase":
			globals.Config.AccountDatabase = index[1]
			break
		case "PNIDCollection":
			globals.Config.PNIDCollection = index[1]
			break
		case "NexAccountsCollection":
			globals.Config.NexAccountsCollection = index[1]
			break
		case "SplatoonDatabase":
			globals.Config.SplatoonDatabase = index[1]
			break
		case "RoomsCollection":
			globals.Config.RoomsCollection = index[1]
			break
		case "SessionsCollection":
			globals.Config.SessionsCollection = index[1]
			break
		case "UsersCollection":
			globals.Config.UsersCollection = index[1]
			break
		case "RegionsCollection":
			globals.Config.RegionsCollection = index[1]
			break
		case "CommunitiesCollection":
			globals.Config.CommunitiesCollection = index[1]
			break
		case "SubscriptionsCollection":
			globals.Config.SubscriptionsCollection = index[1]
			break
		default:
			break
		}
	}
	return config, nil
}
