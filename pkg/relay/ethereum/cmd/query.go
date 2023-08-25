package cmd

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger-labs/yui-relayer/config"
	"github.com/spf13/cobra"

	"github.com/neoiss/ethereum-ibc-relay-chain/pkg/relay/ethereum"
)

// queryCmd represents the chain command
func queryCmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query Commands",
		Long:  "Commands to query useful data on configured chains.",
	}

	cmd.AddCommand(
		queryBalanceCmd(ctx),
	)
	return cmd
}

func queryBalanceCmd(ctx *config.Context) *cobra.Command {
	c := &cobra.Command{
		Use:   "balance [chain-id]",
		Short: "query balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := ctx.Config.GetChain(args[0])
			if err != nil {
				return err
			}
			chain, ok := c.Chain.(*ethereum.Chain)
			if !ok {
				return errors.New("invalid chain-id")
			}
			owner, err := cmd.Flags().GetString(flagOwner)
			if err != nil {
				return err
			}
			bankId, err := cmd.Flags().GetString(flagBankId)
			if err != nil {
				return err
			}
			balance, err := chain.QueryBankBalance(common.HexToAddress(owner), bankId)
			if err != nil {
				return err
			}
			fmt.Printf("%d %s\n", balance, bankId)
			return nil
		},
	}

	c = ownerFlags(c)
	c = bankIdFlags(c)
	return c
}
