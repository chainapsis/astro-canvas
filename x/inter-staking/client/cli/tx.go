package cli

import (
	"bufio"
	"github.com/chainapsis/astro-canvas/x/inter-staking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	interStakingTxCmd := &cobra.Command{
		Use:                types.ModuleName,
		DisableFlagParsing: true,
		RunE:               client.ValidateCmd,
	}

	interStakingTxCmd.AddCommand(flags.PostCommands(GetRegisterCmd(cdc), GetDelegateCmd(cdc))...)

	return interStakingTxCmd
}

func GetRegisterCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "register",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			txBldr, msg, err := BuildRegisterMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsIASourcePort)
	cmd.Flags().AddFlagSet(fsIASourceChannel)

	_ = cmd.MarkFlagRequired(FlagIASourcePort)
	_ = cmd.MarkFlagRequired(FlagIASourceChannel)

	return cmd
}

func GetDelegateCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "delegate [validator-addr] [amount]",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			txBldr, msg, err := BuildDelegateMsg(cliCtx, txBldr, args)
			if err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsTransferSourcePort)
	cmd.Flags().AddFlagSet(fsTransferSourceChannel)
	cmd.Flags().AddFlagSet(fsIASourcePort)
	cmd.Flags().AddFlagSet(fsIASourceChannel)
	cmd.Flags().AddFlagSet(fsCounterpartyBech32Addr)

	_ = cmd.MarkFlagRequired(FlagTransferSourcePort)
	_ = cmd.MarkFlagRequired(FlagTransferSourceChannel)
	_ = cmd.MarkFlagRequired(FlagIASourcePort)
	_ = cmd.MarkFlagRequired(FlagIASourceChannel)
	_ = cmd.MarkFlagRequired(FlagCounterpartyBech32Addr)

	return cmd
}

func BuildRegisterMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	sender := cliCtx.GetFromAddress()

	sourcePort := viper.GetString(FlagIASourcePort)
	sourceChannel := viper.GetString(FlagIASourceChannel)

	msg := types.NewMsgRegister(sourcePort, sourceChannel, sender)

	return txBldr, msg, nil
}

func BuildDelegateMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder, args []string) (auth.TxBuilder, sdk.Msg, error) {
	transferSourcePort := viper.GetString(FlagTransferSourcePort)
	transferSourceChannel := viper.GetString(FlagTransferSourceChannel)
	iaSourcePort := viper.GetString(FlagIASourcePort)
	iaSourceChannel := viper.GetString(FlagIASourceChannel)

	counterpartyBech32Addr := viper.GetString(FlagCounterpartyBech32Addr)

	amount, err := sdk.ParseCoin(args[1])
	if err != nil {
		return txBldr, nil, err
	}

	delAddr := cliCtx.GetFromAddress()
	valAddr, err := sdk.ValAddressFromBech32(args[0])
	if err != nil {
		return txBldr, nil, err
	}

	msg := types.NewMsgDelegate(transferSourcePort, transferSourceChannel, iaSourcePort, iaSourceChannel, counterpartyBech32Addr, delAddr, valAddr, amount)

	return txBldr, msg, nil
}
