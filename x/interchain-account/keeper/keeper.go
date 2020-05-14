package keeper

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/chainapsis/astro-canvas/x/interchain-account/types"
)

const (
	// DefaultPacketTimeoutHeight is the default packet timeout height relative
	// to the current block height. The timeout is disabled when set to 0.
	DefaultPacketTimeoutHeight = 1000 // NOTE: in blocks

	// DefaultPacketTimeoutTimestamp is the default packet timeout timestamp relative
	// to the current block timestamp. The timeout is disabled when set to 0.
	DefaultPacketTimeoutTimestamp = 0 // NOTE: in nanoseconds
)

// Keeper defines the IBC transfer keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.Marshaler

	txCdc codec.Codec
	// This field us used to marshal transaction for counterparty chain.
	// Currently, we support only one counterparty chain per interchain account keeper.
	// TODO: support multiple counterparty codec.
	counterpartyTxCdc codec.Codec

	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	accountKeeper types.AccountKeeper

	scopedKeeper capability.ScopedKeeper

	router types.Router
}

// NewKeeper creates a new IBC transfer Keeper instance
func NewKeeper(
	cdc codec.Marshaler, key sdk.StoreKey,
	txCdc, counterpartyTxCdc codec.Codec, channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper,
	accountKeeper types.AccountKeeper, scopedKeeper capability.ScopedKeeper, router types.Router,
) Keeper {
	return Keeper{
		storeKey:          key,
		cdc:               cdc,
		txCdc:             txCdc,
		counterpartyTxCdc: counterpartyTxCdc,
		channelKeeper:     channelKeeper,
		portKeeper:        portKeeper,
		accountKeeper:     accountKeeper,
		scopedKeeper:      scopedKeeper,
		router:            router,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s/%s", ibctypes.ModuleName, types.ModuleName))
}

func (k Keeper) PacketExecuted(ctx sdk.Context, packet channelexported.PacketI, acknowledgement []byte) error {
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(packet.GetDestPort(), packet.GetDestChannel()))
	if !ok {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, "channel capability could not be retrieved for packet")
	}
	return k.channelKeeper.PacketExecuted(ctx, chanCap, packet, acknowledgement)
}

// IsBound checks if the interchain account module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	// Set the portID into our store so we can retrieve it later
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.PortKey), []byte(portID))

	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, ibctypes.PortPath(portID))
}

// GetPort returns the portID for the transfer module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get([]byte(types.PortKey)))
}

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capability.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}
