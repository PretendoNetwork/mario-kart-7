package globals

import (
	"github.com/PretendoNetwork/mario-kart-7-secure/types"
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/plogger-go"
)

var Logger = plogger.NewLogger()
var NEXServer *nex.Server
var Config *types.ServerConfig

var RegionList = []string{"Worldwide", "Japan", "United States", "Europe", "Korea", "China", "Taiwan"}
var GameModes = []string{"Turf War", "Unk1", "Unk2", "Private Battle", "Unk4"}
var CCList = []string{"Unk", "200cc", "50cc", "100cc", "150cc", "Mirror", "BattleCC"}
var ItemModes = []string{"Unk1", "Unk2", "Unk3", "Unk4", "Unk5", "Normal", "Unk7", "All Items", "Shells Only", "Bananas Only", "Mushrooms Only", "Bob-ombs Only", "No Items", "No Items or Coins", "Frantic"}
var VehicleModes = []string{"All Vehicles", "Karts Only", "Bikes Only"}
var ControllerModes = []string{"Unk", "Tilt Only", "All Controls"}
var DLCModes = []string{"No DLC", "DLC Pack 1 Only", "DLC Pack 2 Only", "Both DLC Packs"}

var GlobalGIDstring = ""
var GlobalHostString = ""
