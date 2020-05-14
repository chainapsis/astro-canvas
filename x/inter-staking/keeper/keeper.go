package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"math"
	"math/rand"
	"strconv"
	"strings"

	interchainaccount "github.com/chainapsis/astro-canvas/x/interchain-account"
	iatypes "github.com/chainapsis/astro-canvas/x/interchain-account/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer"
	transfertypes "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/chainapsis/astro-canvas/x/inter-staking/types"

	"github.com/tendermint/tendermint/libs/bech32"
)

// Inter staking keeper enables users to stake the transfered token to source chain.
// The ideal approach is
// -> (if interchain account doesn't exist on sourch chain, make it)
// -> transfer the token to interchain account on source chain
// -> request the staking msg vai interchain account module if prior process succeeds
// However, it seems that current implementation of relayer by iqlusion doesn't relay the acknowledge packet well
// and can't afford to implement it. So, this assumes that the packet via IBC must succeed
// and doesn't handle the case of failed packet.
type Keeper struct {
	cdc      codec.Marshaler
	storeKey sdk.StoreKey

	// ICS-20 transfer keeper
	transferKeeper          transfer.Keeper
	interchainAccountKeeper interchainaccount.Keeper

	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
}

func NewInterStakingKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey, transferKeeper transfer.Keeper, interchainAccountKeeper interchainaccount.Keeper,
	accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper) Keeper {
	// ensure ibc transfer module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the inter staking module account has not been set")
	}

	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,

		transferKeeper:          transferKeeper,
		interchainAccountKeeper: interchainAccountKeeper,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

// Register interchain account
// This is not needed in ideal process.
// But, currently, this ignores the acknowledge packet
// because current implementation of relayer doesn't relay the acknowledge packet well.
// This ignores the acknowledgement of packet, and assumes the packets must succeed.
// But, it is hard to assume that registering interchain account must succeed.
// So, ask the user to register an interchain account at first.
func (keeper Keeper) RegisterInterchainAccount(ctx sdk.Context, sender sdk.AccAddress, sourcePort, sourceChannel string) error {
	// Use the block height as random seed to produce random number.
	r := rand.New(rand.NewSource(ctx.BlockHeight()))

	randInt := r.Int63()
	salt := strconv.FormatInt(randInt, 10)
	err := keeper.interchainAccountKeeper.CreateInterchainAccount(ctx, sourcePort, sourceChannel, salt)
	if err != nil {
		return err
	}

	address, err := keeper.interchainAccountKeeper.GenerateAddress(iatypes.GetIdentifier(sourcePort, sourceChannel), salt)
	if err != nil {
		return err
	}

	kvStore := ctx.KVStore(keeper.storeKey)
	prefixStore := prefix.NewStore(kvStore, []byte("ia/"))

	key := []byte(fmt.Sprintf("%s/%s/%s", sourcePort, sourceChannel, sender.String()))
	if prefixStore.Has(key) {
		return types.ErrIAAccountAlreadyExist
	}
	prefixStore.Set(key, address)

	ctx.EventManager().EmitEvent(sdk.NewEvent("register-interchain-account", sdk.NewAttribute("expected-address", sdk.AccAddress(address).String())))

	return nil
}

func (keeper Keeper) GetInterchainAccount(ctx sdk.Context, address sdk.AccAddress, sourcePort, sourceChannel string) ([]byte, error) {
	kvStore := ctx.KVStore(keeper.storeKey)
	prefixStore := prefix.NewStore(kvStore, []byte("ia/"))

	key := []byte(fmt.Sprintf("%s/%s/%s", sourcePort, sourceChannel, address.String()))
	if !prefixStore.Has(key) {
		return []byte{}, types.ErrIAAccountNotExist
	}
	bz := prefixStore.Get(key)

	return bz, nil
}

// Perform delegate to other chain.
// But, packets must be executed in order of send packet -> delegate packet.
// Currently, this ingores constraint due to the lack of working time.
// Relayer must handle this order.
func (keeper Keeper) Delegate(ctx sdk.Context, counterpartyBech32Addr string, delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin,
	transferSourcePort, transferSourceChannel, iaSourcePort, iaSourceChannel string) error {
	iaAccount, err := keeper.GetInterchainAccount(ctx, delAddr, iaSourcePort, iaSourceChannel)
	if err != nil {
		return err
	}

	recipient, err := bech32.ConvertAndEncode(counterpartyBech32Addr, iaAccount)
	if err != nil {
		return err
	}

	err = keeper.transferKeeper.SendTransfer(ctx, transferSourcePort, transferSourceChannel, math.MaxUint64, sdk.Coins{amount}, delAddr, recipient)
	if err != nil {
		return err
	}

	prefixDenom := transfertypes.GetDenomPrefix(transferSourcePort, transferSourceChannel)
	amount = sdk.NewCoin(strings.Replace(amount.Denom, prefixDenom, "", 1), amount.Amount)

	err = keeper.interchainAccountKeeper.RequestRunTx(ctx, iaSourcePort, iaSourceChannel, iatypes.CosmosSdkChainType, staking.NewMsgDelegate(delAddr, valAddr, amount))
	if err != nil {
		return err
	}

	mint := sdk.Coins{sdk.NewCoin(fmt.Sprintf("%s/%s/stake", iaSourcePort, iaSourceChannel), amount.Amount)}
	// Mint coin as the proof of staking
	err = keeper.bankKeeper.MintCoins(ctx, types.ModuleName, mint)
	if err != nil {
		return err
	}

	// TODO: Send the minted token to sender when the packet is exeucted successfully.
	err = keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delAddr, mint)
	if err != nil {
		return err
	}

	return nil
}
