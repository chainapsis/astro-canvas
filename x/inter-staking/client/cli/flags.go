package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagTransferSourcePort    = "transfer-source-port"
	FlagTransferSourceChannel = "transfer-source-channel"
	FlagIASourcePort          = "ia-source-port"
	FlagIASourceChannel       = "ia-source-channel"

	FlagCounterpartyBech32Addr = "counterparty-bech32-addr"
)

// common flagsets to add to various functions
var (
	fsTransferSourcePort    = flag.NewFlagSet("", flag.ContinueOnError)
	fsTransferSourceChannel = flag.NewFlagSet("", flag.ContinueOnError)
	fsIASourcePort          = flag.NewFlagSet("", flag.ContinueOnError)
	fsIASourceChannel       = flag.NewFlagSet("", flag.ContinueOnError)

	fsCounterpartyBech32Addr = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsTransferSourcePort.String(FlagTransferSourcePort, "", "Source port for ics-20 transfer")
	fsTransferSourceChannel.String(FlagTransferSourceChannel, "", "Source channel for ics-20 transfer")
	fsIASourcePort.String(FlagIASourcePort, "", "Source port for interchain account")
	fsIASourceChannel.String(FlagIASourceChannel, "", "Source channel for interchain account")

	fsCounterpartyBech32Addr.String(FlagCounterpartyBech32Addr, "", "The Bech32 addr of counterparty chain")
}
