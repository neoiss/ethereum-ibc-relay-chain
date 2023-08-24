package cmd

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hyperledger-labs/yui-relayer/config"
	"github.com/hyperledger-labs/yui-relayer/core"
	"github.com/spf13/cobra"

	"github.com/datachainlab/ethereum-ibc-relay-chain/pkg/relay/ethereum"
)

// txCmd represents the chain command
func txCmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Tx Commands",
		Long:  "Commands to send tx.",
	}

	cmd.AddCommand(
		xfersend(ctx),
	)
	return cmd
}

// rly harmony tx transfer ibc01 ibc1 --amount 100 --denom ${HMY_TOKEN_DENOM} --receiver ${TM_ADDRESS}
func xfersend(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [path-name] [chain-id]",
		Short: "Initiate a transfer from one chain to another",
		Long: "Sends the first step to transfer tokens in an IBC transfer." +
			" The created packet must be relayed to another chain",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := ctx.Config.Paths.Get(args[0])
			if err != nil {
				return err
			}
			c, err := ctx.Config.GetChain(args[1])
			if err != nil {
				return err
			}
			_, ok := c.Chain.(*ethereum.Chain)
			if !ok {
				return errors.New("invalid chain-id")
			}

			amount, err := cmd.Flags().GetUint64(flagAmount)
			if err != nil {
				return err
			}
			// XXX want to support all denom format
			d, err := cmd.Flags().GetString(flagDenom)
			if err != nil {
				return err
			}
			denom := transfertypes.ParseDenomTrace(d)
			token := sdk.Coin{
				Denom:  denom.GetFullDenomPath(),
				Amount: sdk.Int(sdk.NewUint(amount)),
			}
			if denom.Path != "" {
				token.Denom = denom.IBCDenom()
			}

			// Bech32 address string
			receiver, err := cmd.Flags().GetString(flagReceiver)
			if err != nil {
				return err
			}
			receiverAcc, err := sdk.AccAddressFromBech32(receiver)
			if err != nil {
				return err
			}
			fmt.Printf("receiver: %s", hexutil.Encode(receiverAcc.Bytes()))

			memo, err := cmd.Flags().GetString(flagMemo)
			if err != nil {
				return err
			}

			tx := core.RelayMsgs{
				Src: []sdk.Msg{},
				Dst: []sdk.Msg{
					transfertypes.NewMsgTransfer(
						path.Dst.PortID,
						path.Dst.ChannelID,
						token,
						"", // not used
						hexutil.Encode(receiverAcc.Bytes()),
						// TODO timeout height
						types.NewHeight(0, 10000),
						// TODO timeout timestamp
						uint64(0),
						memo,
					),
				},
			}

			if tx.Send(nil, c); !tx.Succeeded {
				return fmt.Errorf("failed to send transfer message")
			}
			return nil
		},
	}
	return sendTransferFlags(cmd)
}
