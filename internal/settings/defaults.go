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

	// DefaultLogTarget is the default output location of log output.
	defaultRegisterURL = "http://localhost:9020"
)
