package cli

import (
	"bufio"
	"github.com/chainapsis/astro-canvas/x/canvas/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	canvasTxCmd := &cobra.Command{
		Use:                types.ModuleName,
		DisableFlagParsing: true,
		RunE:               client.ValidateCmd,
	}

	canvasTxCmd.AddCommand(flags.PostCommands(GetCreateCanvasCmd(cdc), GetPaintCmd(cdc))...)

	return canvasTxCmd
}

func GetCreateCanvasCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create-canvas [id] [width] [height] [refundDuration] [denomPrefix] [price]",
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			txBldr, msg, err := BuildCreateCanvasMsg(cliCtx, txBldr, args)
			if err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

func GetPaintCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "paint [id] [x] [y] [amount]",
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			txBldr, msg, err := BuildPaintMsg(cliCtx, txBldr, args)
			if err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}

func BuildCreateCanvasMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder, args []string) (auth.TxBuilder, sdk.Msg, error) {
	//"create-canvas [id] [width] [height] [refundDuration] [denomPrefix] [price]",

	sender := cliCtx.GetFromAddress()

	id := args[0]
	width, err := strconv.Atoi(args[1])
	if err != nil {
		return txBldr, nil, err
	}
	height, err := strconv.Atoi(args[2])
	if err != nil {
		return txBldr, nil, err
	}

	refundDuration, err := time.ParseDuration(args[3])
	if err != nil {
		return txBldr, nil, err
	}

	denomPrefix := args[4]

	price, err := strconv.Atoi(args[5])
	if err != nil {
		return txBldr, nil, err
	}

	msg := types.NewMsgCreateCanvas(id, uint64(width), uint64(height), refundDuration, denomPrefix, uint64(price), sender)

	return txBldr, msg, nil
}

func BuildPaintMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder, args []string) (auth.TxBuilder, sdk.Msg, error) {
	id := args[0]

	x, err := strconv.Atoi(args[1])
	if err != nil {
		return txBldr, nil, err
	}

	y, err := strconv.Atoi(args[2])
	if err != nil {
		return txBldr, nil, err
	}

	amount, err := sdk.ParseCoin(args[3])
	if err != nil {
		return txBldr, nil, err
	}

	sender := cliCtx.GetFromAddress()

	msg := types.NewMsgPaint(id, uint64(x), uint64(y), amount, sender)

	return txBldr, msg, nil
}
