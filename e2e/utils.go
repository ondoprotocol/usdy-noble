package e2e

import (
	"context"
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/jsonpb"
	_ "github.com/noble-assets/aura/x/aura"
	"github.com/noble-assets/aura/x/aura/types"
	"github.com/noble-assets/aura/x/aura/types/blocklist"
	"github.com/strangelove-ventures/interchaintest/v4"
	"github.com/strangelove-ventures/interchaintest/v4/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v4/ibc"
	"github.com/strangelove-ventures/interchaintest/v4/relayer/rly"
	"github.com/strangelove-ventures/interchaintest/v4/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

type Wrapper struct {
	chain *cosmos.CosmosChain
	gaia  *cosmos.CosmosChain

	owner   ibc.Wallet
	minter  ibc.Wallet
	burner  ibc.Wallet
	alice   ibc.Wallet
	bob     ibc.Wallet
	charlie ibc.Wallet
}

func Suite(t *testing.T, wrapper *Wrapper, ibcEnabled bool) (ctx context.Context) {
	ctx = context.Background()
	logger := zaptest.NewLogger(t)
	reporter := testreporter.NewNopReporter()
	execReporter := reporter.RelayerExecReporter(t)
	client, network := interchaintest.DockerSetup(t)

	numValidators, numFullNodes := 1, 0

	specs := []*interchaintest.ChainSpec{
		{
			Name:          "aura",
			Version:       "local",
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
			ChainConfig: ibc.ChainConfig{
				Type:    "cosmos",
				Name:    "aura",
				ChainID: "aura-1",
				Images: []ibc.DockerImage{
					{
						Repository: "aura-simd",
						Version:    "local",
						UidGid:     "1025:1025",
					},
				},
				Bin:            "simd",
				Bech32Prefix:   "noble",
				Denom:          "ustake",
				GasPrices:      "0.0ustake",
				GasAdjustment:  5,
				TrustingPeriod: "504h",
				NoHostMount:    false,
				PreGenesis: func(cfg ibc.ChainConfig) (err error) {
					validator := wrapper.chain.Validators[0]
					ONE := sdk.NewCoins(sdk.NewInt64Coin("ustake", 1_000_000))

					wrapper.owner, err = wrapper.chain.BuildWallet(ctx, "owner", "")
					if err != nil {
						return err
					}
					err = validator.AddGenesisAccount(ctx, wrapper.owner.FormattedAddress(), ONE)
					if err != nil {
						return err
					}

					wrapper.minter, err = wrapper.chain.BuildWallet(ctx, "minter", "")
					if err != nil {
						return err
					}
					err = validator.AddGenesisAccount(ctx, wrapper.minter.FormattedAddress(), ONE)
					if err != nil {
						return err
					}

					wrapper.burner, err = wrapper.chain.BuildWallet(ctx, "burner", "")
					if err != nil {
						return err
					}
					err = validator.AddGenesisAccount(ctx, wrapper.burner.FormattedAddress(), ONE)
					if err != nil {
						return err
					}

					wrapper.alice, err = wrapper.chain.BuildWallet(ctx, "alice", "")
					if err != nil {
						return err
					}
					err = validator.AddGenesisAccount(ctx, wrapper.alice.FormattedAddress(), ONE)
					if err != nil {
						return err
					}

					wrapper.bob, err = wrapper.chain.BuildWallet(ctx, "bob", "")
					if err != nil {
						return err
					}
					err = validator.AddGenesisAccount(ctx, wrapper.bob.FormattedAddress(), ONE)
					if err != nil {
						return err
					}

					return nil
				},
				ModifyGenesis: func(cfg ibc.ChainConfig, bz []byte) ([]byte, error) {
					ONE := sdk.NewInt(1_000_000_000_000_000_000)

					changes := []cosmos.GenesisKV{
						{Key: "app_state.aura.blocklist_state.owner", Value: wrapper.owner.FormattedAddress()},
						{Key: "app_state.aura.owner", Value: wrapper.owner.FormattedAddress()},
						{Key: "app_state.aura.minters", Value: []types.Minter{
							{Address: wrapper.minter.FormattedAddress(), Allowance: ONE},
						}},
						{Key: "app_state.aura.burners", Value: []types.Burner{
							{Address: wrapper.burner.FormattedAddress(), Allowance: ONE},
						}},
					}

					return cosmos.ModifyGenesis(changes)(cfg, bz)
				},
				ModifyGenesisAmounts: func() (sdk.Coin, sdk.Coin) {
					ONE := sdk.NewInt64Coin("ustake", 1_000_000)
					return ONE, ONE
				},
			},
		},
	}
	if ibcEnabled {
		specs = append(specs, &interchaintest.ChainSpec{
			Name:          "ibc-go-simd",
			Version:       "v4.5.0",
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
			ChainConfig: ibc.ChainConfig{
				PreGenesis: func(cfg ibc.ChainConfig) (err error) {
					validator := wrapper.gaia.Validators[0]
					ONE := sdk.NewCoins(sdk.NewInt64Coin(cfg.Denom, 1_000_000))

					wrapper.charlie, err = wrapper.gaia.BuildWallet(ctx, "owner", "")
					if err != nil {
						return err
					}
					err = validator.AddGenesisAccount(ctx, wrapper.charlie.FormattedAddress(), ONE)
					if err != nil {
						return err
					}

					return nil
				},
			},
		})
	}
	factory := interchaintest.NewBuiltinChainFactory(logger, specs)

	chains, err := factory.Chains(t.Name())
	require.NoError(t, err)

	noble := chains[0].(*cosmos.CosmosChain)
	wrapper.chain = noble

	interchain := interchaintest.NewInterchain().AddChain(noble)
	var relayer *rly.CosmosRelayer
	if ibcEnabled {
		relayer = interchaintest.NewBuiltinRelayerFactory(
			ibc.CosmosRly,
			logger,
		).Build(t, client, network).(*rly.CosmosRelayer)

		gaia := chains[1].(*cosmos.CosmosChain)
		wrapper.gaia = gaia

		interchain = interchain.
			AddChain(gaia).
			AddRelayer(relayer, "relayer").
			AddLink(interchaintest.InterchainLink{
				Chain1:  noble,
				Chain2:  gaia,
				Relayer: relayer,
				Path:    "transfer",
			})
	}
	require.NoError(t, interchain.Build(ctx, execReporter, interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
	}))

	t.Cleanup(func() {
		_ = interchain.Close()
	})

	if ibcEnabled {
		require.NoError(t, relayer.StartRelayer(ctx, execReporter))
	}

	return
}

//

func EnsureBlocked(t *testing.T, wrapper Wrapper, ctx context.Context, address string) {
	validator := wrapper.chain.Validators[0]

	raw, _, err := validator.ExecQuery(ctx, "aura", "blocklist", "address", address)
	require.NoError(t, err)

	var res blocklist.QueryAddressResponse
	require.NoError(t, json.Unmarshal(raw, &res))

	require.True(t, res.Blocked)
}

func EnsureBurner(t *testing.T, wrapper Wrapper, ctx context.Context, address string, allowance sdk.Int) {
	validator := wrapper.chain.Validators[0]

	raw, _, err := validator.ExecQuery(ctx, "aura", "burners")
	require.NoError(t, err)

	var res types.QueryBurnersResponse
	require.NoError(t, jsonpb.UnmarshalString(string(raw), &res))

	require.Contains(t, res.Burners, types.Burner{Address: address, Allowance: allowance})
}

func EnsureMinter(t *testing.T, wrapper Wrapper, ctx context.Context, address string, allowance sdk.Int) {
	validator := wrapper.chain.Validators[0]

	raw, _, err := validator.ExecQuery(ctx, "aura", "minters")
	require.NoError(t, err)

	var res types.QueryMintersResponse
	require.NoError(t, jsonpb.UnmarshalString(string(raw), &res))

	require.Contains(t, res.Minters, types.Minter{Address: address, Allowance: allowance})
}

func EnsureBlockedChannel(t *testing.T, wrapper Wrapper, ctx context.Context, channel string) {
	validator := wrapper.chain.Validators[0]

	raw, _, err := validator.ExecQuery(ctx, "aura", "blocked-channels")
	require.NoError(t, err)

	var res types.QueryBlockedChannelsResponse
	require.NoError(t, json.Unmarshal(raw, &res))

	require.Contains(t, res.BlockedChannels, channel)
}
