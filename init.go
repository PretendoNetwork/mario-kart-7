package main

//"io/ioutil"

var hmacSecret []byte

func init() {
	//var err error

	/*hmacSecret, err = ioutil.ReadFile("secret.key")
	if err != nil {
		panic(err)
	}*/

	config, _ = ImportConfigFromFile("secure.config")
	connectMongo()
}
