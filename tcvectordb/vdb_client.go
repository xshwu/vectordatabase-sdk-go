package tcvectordb

// VDBClient tencent vectordb client
type VDBCLient struct {
	DatabaseInterface
	FlatInterface
	FlatIndexInterface

	cli SdkClient
}

// NewVDBClient new VDBClient with external SdkClient implement
func NewVDBClient(cli SdkClient) *VDBCLient {

	databaseImpl := new(implementerDatabase)
	databaseImpl.SdkClient = cli

	flatImpl := new(implementerFlatDocument)
	flatImpl.SdkClient = cli

	flatIndexImpl := new(implementerFlatIndex)
	flatIndexImpl.SdkClient = cli

	return &VDBCLient{
		cli: cli,

		DatabaseInterface:  databaseImpl,
		FlatInterface:      flatImpl,
		FlatIndexInterface: flatIndexImpl,
	}
}
