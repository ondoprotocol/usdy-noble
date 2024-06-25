package cli

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/aura/x/aura/types"
	"github.com/spf13/cobra"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Transactions commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetBlocklistTxCmd())

	cmd.AddCommand(TxBurn())
	cmd.AddCommand(TxMint())
	cmd.AddCommand(TxPause())
	cmd.AddCommand(TxUnpause())
	cmd.AddCommand(TxTransferOwnership())
	cmd.AddCommand(TxAcceptOwnership())
	cmd.AddCommand(TxAddBurner())
	cmd.AddCommand(TxRemoveBurner())
	cmd.AddCommand(TxSetBurnerAllowance())
	cmd.AddCommand(TxAddMinter())
	cmd.AddCommand(TxRemoveMinter())
	cmd.AddCommand(TxSetMinterAllowance())
	cmd.AddCommand(TxAddPauser())
	cmd.AddCommand(TxRemovePauser())

	return cmd
}

func TxBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [from] [amount]",
		Short: "Transaction that burns tokens",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := &types.MsgBurn{
				Signer: clientCtx.GetFromAddress().String(),
				From:   args[0],
				Amount: amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [to] [amount]",
		Short: "Transaction that mints tokens",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := &types.MsgMint{
				Signer: clientCtx.GetFromAddress().String(),
				To:     args[0],
				Amount: amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxPause() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause",
		Short: "Transaction that pauses the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgPause{
				Signer: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxUnpause() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unpause",
		Short: "Transaction that unpauses the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgUnpause{
				Signer: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxTransferOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ownership [new-owner]",
		Short: "Transfer ownership of module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgTransferOwnership{
				Signer:   clientCtx.GetFromAddress().String(),
				NewOwner: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxAcceptOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept-ownership",
		Short: "Accept ownership of module",
		Long:  "Accept ownership of module, assuming there is an pending ownership transfer",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgAcceptOwnership{
				Signer: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxAddBurner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-burner [burner] [allowance]",
		Short: "Add a new burner with an initial allowance",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allowance, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid allowance")
			}

			msg := &types.MsgAddBurner{
				Signer:    clientCtx.GetFromAddress().String(),
				Burner:    args[0],
				Allowance: allowance,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxRemoveBurner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-burner [burner]",
		Short: "Removes an existing burner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRemoveBurner{
				Signer: clientCtx.GetFromAddress().String(),
				Burner: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxSetBurnerAllowance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-burner-allowance [burner] [allowance]",
		Short: "Sets an existing burner's allowance",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allowance, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid allowance")
			}

			msg := &types.MsgSetBurnerAllowance{
				Signer:    clientCtx.GetFromAddress().String(),
				Burner:    args[0],
				Allowance: allowance,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxAddMinter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-minter [minter] [allowance]",
		Short: "Add a new minter with an initial allowance",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allowance, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid allowance")
			}

			msg := &types.MsgAddMinter{
				Signer:    clientCtx.GetFromAddress().String(),
				Minter:    args[0],
				Allowance: allowance,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxRemoveMinter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-minter [minter]",
		Short: "Removes an existing minter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRemoveMinter{
				Signer: clientCtx.GetFromAddress().String(),
				Minter: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxSetMinterAllowance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-minter-allowance [minter] [allowance]",
		Short: "Sets an existing minter's allowance",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			allowance, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid allowance")
			}

			msg := &types.MsgSetMinterAllowance{
				Signer:    clientCtx.GetFromAddress().String(),
				Minter:    args[0],
				Allowance: allowance,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxAddPauser() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-pauser [pauser]",
		Short: "Add a new pauser",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgAddPauser{
				Signer: clientCtx.GetFromAddress().String(),
				Pauser: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxRemovePauser() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-pauser [pauser]",
		Short: "Removes an existing pauser",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRemovePauser{
				Signer: clientCtx.GetFromAddress().String(),
				Pauser: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
