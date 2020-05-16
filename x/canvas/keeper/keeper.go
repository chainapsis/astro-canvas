package keeper

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/chainapsis/astro-canvas/x/canvas/types"
)

type Keeper struct {
	cdc      codec.Marshaler
	storeKey sdk.StoreKey

	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
}

func NewCanvasKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper) Keeper {
	// ensure canvas module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the canvas module account has not been set")
	}

	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

func (keeper Keeper) CreateCanvas(ctx sdk.Context, id string, width uint64, height uint64, refundDuration time.Duration, allowDenomPrefix string, priceForPoint uint64) error {
	kvStore := ctx.KVStore(keeper.storeKey)
	prefixStore := prefix.NewStore(kvStore, []byte("canvas/"))

	if prefixStore.Has([]byte(id)) {
		return types.ErrCanvasAlreadyExist
	}

	canvas := types.Canvas{
		Id:               id,
		Width:            width,
		Height:           height,
		RefundDuration:   refundDuration,
		AllowDenomPrefix: allowDenomPrefix,
		PriceForPoint:    priceForPoint,
	}

	bz, err := keeper.cdc.MarshalBinaryBare(&canvas)
	if err != nil {
		return err
	}

	prefixStore.Set([]byte(id), bz)
	return nil
}

func (keeper Keeper) GetCanvas(ctx sdk.Context, id string) (types.Canvas, error) {
	kvStore := ctx.KVStore(keeper.storeKey)
	prefixStore := prefix.NewStore(kvStore, []byte("canvas/"))

	if !prefixStore.Has([]byte(id)) {
		return types.Canvas{}, types.ErrCanvasNotExist
	}

	bz := prefixStore.Get([]byte(id))

	canvas := types.Canvas{}
	err := keeper.cdc.UnmarshalBinaryBare(bz, &canvas)
	if err != nil {
		return types.Canvas{}, err
	}

	return canvas, nil
}

func (keeper Keeper) Paint(ctx sdk.Context, id string, x uint64, y uint64, amount sdk.Coin, sender sdk.AccAddress) error {
	canvas, err := keeper.GetCanvas(ctx, id)
	if err != nil {
		return err
	}

	if x >= canvas.Width || y >= canvas.Height {
		return types.ErrPointGetOut
	}

	if !strings.HasPrefix(amount.Denom, canvas.AllowDenomPrefix) {
		return types.ErrInvalidDenomPrefix
	}

	if amount.Amount.IsNegative() || !amount.Amount.Equal(sdk.NewInt(int64(canvas.PriceForPoint))) {
		return types.ErrInvalidAmount
	}

	refundTime := ctx.BlockTime().Add(canvas.RefundDuration)

	refundData := types.RefundData{
		Recipient:  sender,
		Amount:     amount,
		RefundTime: refundTime.Unix(),
	}

	bz, err := keeper.cdc.MarshalBinaryBare(&refundData)
	if err != nil {
		return err
	}

	kvStore := ctx.KVStore(keeper.storeKey)
	prefixStore := prefix.NewStore(kvStore, []byte("refund/"))

	refundDataKey := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(refundDataKey, keeper.getIncrementalSequence(ctx))
	prefixStore.Set(refundDataKey, bz)

	err = keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{amount})
	if err != nil {
		return err
	}

	prefixStore = prefix.NewStore(kvStore, []byte(fmt.Sprintf("point/%s/", id)))

	point := types.Point{
		X:     x,
		Y:     y,
		Color: amount.Denom,
	}

	bz, err = keeper.cdc.MarshalBinaryBare(&point)
	if err != nil {
		return err
	}

	prefixStore.Set([]byte(fmt.Sprintf("%d/%d", x, y)), bz)
	return nil
}

func (keeper Keeper) GetPoints(ctx sdk.Context, id string) ([]types.Point, error) {
	_, err := keeper.GetCanvas(ctx, id)
	if err != nil {
		return nil, err
	}

	points := make([]types.Point, 0)

	kvStore := ctx.KVStore(keeper.storeKey)
	prefixStore := prefix.NewStore(kvStore, []byte(fmt.Sprintf("point/%s/", id)))

	iter := prefixStore.Iterator([]byte{}, []byte{255})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		point := types.Point{}
		err := keeper.cdc.UnmarshalBinaryBare(iter.Value(), &point)
		if err != nil {
			return nil, err
		}

		points = append(points, point)
	}

	return points, nil
}

// Should be executed per each BeginBlock
func (keeper Keeper) Refund(ctx sdk.Context) error {
	kvStore := ctx.KVStore(keeper.storeKey)
	prefixStore := prefix.NewStore(kvStore, []byte("refund/"))

	refundDatas := make([]struct {
		types.RefundData
		Key []byte
	}, 0)

	iter := prefixStore.Iterator([]byte{}, []byte{255})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		refundData := types.RefundData{}
		err := keeper.cdc.UnmarshalBinaryBare(iter.Value(), &refundData)
		if err != nil {
			return err
		}

		refundDatas = append(refundDatas, struct {
			types.RefundData
			Key []byte
		}{
			RefundData: refundData,
			Key:        iter.Key(),
		})
	}

	for _, refundData := range refundDatas {
		refundTime := time.Unix(refundData.RefundTime, 0)

		if ctx.BlockTime().After(refundTime) {
			err := keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, refundData.Recipient, sdk.Coins{refundData.Amount})
			if err != nil {
				return err
			}

			prefixStore.Delete(refundData.Key)
		}
	}

	return nil
}

func (keeper Keeper) getIncrementalSequence(ctx sdk.Context) uint64 {
	kvStore := ctx.KVStore(keeper.storeKey)
	prefixStore := prefix.NewStore(kvStore, []byte("seq/"))

	bz := prefixStore.Get([]byte("seq"))
	if len(bz) == 0 {
		buf := make([]byte, binary.MaxVarintLen64)
		binary.PutUvarint(buf, 1)
		prefixStore.Set([]byte("seq"), buf)

		return 0
	}

	seq, _ := binary.Uvarint(bz)

	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, seq+1)
	prefixStore.Set([]byte("seq"), buf)

	return seq
}
