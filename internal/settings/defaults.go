package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Client Settings

const (
	// DefaultEstablishmentTTL is the default Time To Live used with Client - Gateway estalishment messages.
	defaultEstablishmentTTL = int64(100)

	// DefaultLogLevel is the default amount of logging to show.
	defaultLogLevel = "trace"

	// DefaultLogTarget is the default output location of log output.
	defaultLogTarget = "STDOUT"

	// DefaultRegisterURL is the default location of the Register service. 
	// register:9020 is the value that will work for the integration test system.
	defaultRegisterURL = "http://register:9020"
)
