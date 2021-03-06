package utils

import (
	"github.com/urfave/cli"
)

var (
	// ValidatorFlag determines if a node will run as voter. Participant should have deposited 32ETH and ready to perform proposer and attester duties.
	ValidatorFlag = cli.BoolFlag{
		Name:  "voter",
		Usage: "Whether or not to run the node as beacon chain voter",
	}
	// SimulatorFlag determines if a node will run only as a simulator service.
	SimulatorFlag = cli.BoolFlag{
		Name:  "simulator",
		Usage: "Whether or not to run the node as a simple simulator of beacon blocks over p2p",
	}
	// Web3ProviderFlag defines a flag for a mainchain RPC endpoint.
	Web3ProviderFlag = cli.StringFlag{
		Name:  "web3provider",
		Usage: "A mainchain web3 provider string endpoint. Can either be an IPC file string or a WebSocket endpoint. Uses WebSockets by default at ws://127.0.0.1:8546. Cannot be an HTTP endpoint.",
		Value: "ws://127.0.0.1:8546",
	}
	// VrcContractFlag defines a flag for VRC contract address.
	VrcContractFlag = cli.StringFlag{
		Name:  "vrcaddr",
		Usage: "Validator registration contract address. Beacon chain node will listen logs coming from VRC to determine when voter is eligible to participate.",
	}
	// PubKeyFlag defines a flag for voter's public key on the mainchain
	PubKeyFlag = cli.StringFlag{
		Name:  "pubkey",
		Usage: "Validator's public key. Beacon chain node will listen to VRC log to determine when registration has completed based on this public key address.",
	}
	// RPCPort defines a beacon node RPC port to open.
	RPCPort = cli.StringFlag{
		Name:  "rpc-port",
		Usage: "RPC port exposed by a beacon node",
		Value: "4000",
	}
	// CertFlag defines a flag for the node's TLS certificate.
	CertFlag = cli.StringFlag{
		Name:  "tls-cert",
		Usage: "Certificate for secure gRPC. Pass this and the tls-key flag in order to use gRPC securely.",
	}
	// KeyFlag defines a flag for the node's TLS key.
	KeyFlag = cli.StringFlag{
		Name:  "tls-key",
		Usage: "Key for secure gRPC. Pass this and the tls-cert flag in order to use gRPC securely.",
	}
)
