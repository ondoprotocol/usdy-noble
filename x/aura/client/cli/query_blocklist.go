package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ondoprotocol/usdy-noble/v2/x/aura/types/blocklist"
	"github.com/spf13/cobra"
)

func GetBlocklistQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "blocklist",
		Short:                      "Querying commands for the blocklist submodule",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(QueryBlocklistOwner())
	cmd.AddCommand(QueryBlockedAddresses())
	cmd.AddCommand(QueryBlockedAddress())

	return cmd
}

func QueryBlocklistOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner",
		Short: "Query the submodule's owner",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := blocklist.NewQueryClient(clientCtx)

			res, err := queryClient.Owner(context.Background(), &blocklist.QueryOwner{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryBlockedAddresses() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addresses",
		Short: "Query for all blocked addresses",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := blocklist.NewQueryClient(clientCtx)

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Addresses(context.Background(), &blocklist.QueryAddresses{
				Pagination: pagination,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryBlockedAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "address",
		Short: "Query if an address is blocked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := blocklist.NewQueryClient(clientCtx)

			res, err := queryClient.Address(context.Background(), &blocklist.QueryAddress{
				Address: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
