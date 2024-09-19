package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ondoprotocol/usdy-noble/v2/x/aura/types/blocklist"
	"github.com/spf13/cobra"
)

func GetBlocklistTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "blocklist",
		Short:                      "Transactions commands for the blocklist submodule",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(TxBlocklistTransferOwnership())
	cmd.AddCommand(TxBlocklistAcceptOwnership())
	cmd.AddCommand(TxAddToBlocklist())
	cmd.AddCommand(TxRemoveFromBlocklist())

	return cmd
}

func TxBlocklistTransferOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ownership [new-owner]",
		Short: "Transfer ownership of submodule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &blocklist.MsgTransferOwnership{
				Signer:   clientCtx.GetFromAddress().String(),
				NewOwner: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxBlocklistAcceptOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept-ownership",
		Short: "Accept ownership of submodule",
		Long:  "Accept ownership of submodule, assuming there is an pending ownership transfer",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &blocklist.MsgAcceptOwnership{
				Signer: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxAddToBlocklist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-to-blocklist [addresses ...]",
		Short: "Add a list of accounts to the blocklist",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &blocklist.MsgAddToBlocklist{
				Signer:   clientCtx.GetFromAddress().String(),
				Accounts: args,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxRemoveFromBlocklist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-from-blocklist [addresses ...]",
		Short: "Remove a list of accounts from the blocklist",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &blocklist.MsgRemoveFromBlocklist{
				Signer:   clientCtx.GetFromAddress().String(),
				Accounts: args,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
