package cli

const (
	// ServerHost is where the server is running
	ServerHost = "controller.bytescheme.com"
)

var (
	storeCmdParams = &storeCommandParams{
		serverParams: serverParams{
			host:   ServerHost,
			port:   443,
			scheme: "https",
			apiKey: "Abomcha@123",
		},
	}

	serviceCmdParams = &serviceCommandParams{
		serverParams: serverParams{
			host:   ServerHost,
			port:   443,
			scheme: "https",
			apiKey: "Abomcha@123",
		},
	}

	controllerCmdParams = &controllerCommandParams{
		serverParams: serverParams{
			host:   ServerHost,
			port:   443,
			scheme: "https",
			apiKey: "Abomcha@123",
		},
	}
)

type serverParams struct {
	host   string
	port   int
	apiKey string
	scheme string
}

type storeCommandParams struct {
	serverParams
	isLocal  bool
	key      string
	isPrefix bool
	value    string
	file     string
}

type serviceCommandParams struct {
	serverParams
}

type controllerCommandParams struct {
	serverParams
	controllerID string
	pinID        int
	pinHigh      bool
}
