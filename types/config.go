package types

type ServerConfig struct {
	ServerName              string
	ServerPort              string
	PrudpVersion            int
	SignatureVersion        int
	KerberosKeySize         int
	AccessKey               string
	NexVersion              int
	NEXDatabaseIP           string
	NEXDatabasePort         string
	AccountDatabaseIP       string
	AccountDatabasePort     string
	DatabaseUseAuth         bool
	DatabaseUsername        string
	DatabasePassword        string
	AccountDatabase         string
	PNIDCollection          string
	NexAccountsCollection   string
	SplatoonDatabase        string
	RoomsCollection         string
	SessionsCollection      string
	UsersCollection         string
	RegionsCollection       string
	CommunitiesCollection   string
	SubscriptionsCollection string
}
