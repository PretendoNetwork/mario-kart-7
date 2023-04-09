package main

import (
	"github.com/PretendoNetwork/mario-kart-7-secure/database"
	"github.com/PretendoNetwork/mario-kart-7-secure/globals"
	"github.com/PretendoNetwork/mario-kart-7-secure/utility"
)

func init() {

	globals.Config, _ = utility.ImportConfigFromFile("secure.config")

	database.ConnectAll()
}
