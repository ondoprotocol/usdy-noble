package simapp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-go/v4/modules/apps/transfer"
	transferkeeper "github.com/cosmos/ibc-go/v4/modules/apps/transfer/keeper"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v4/modules/core"
	porttypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/noble-assets/aura/x/aura"
	aurakeeper "github.com/noble-assets/aura/x/aura/keeper"
	auratypes "github.com/noble-assets/aura/x/aura/types"
)

var (
	DefaultNodeHome string

	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		genutil.AppModuleBasic{},
		params.AppModuleBasic{},
		staking.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		capability.AppModuleBasic{},
		ibc.AppModuleBasic{},
		transfer.AppModuleBasic{},
		aura.AppModuleBasic{},
	)

	permissions = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		transfertypes.ModuleName:       {authtypes.Burner, authtypes.Minter},
		auratypes.ModuleName:           {authtypes.Burner, authtypes.Minter},
	}
)

var (
	_ simapp.App              = (*SimApp)(nil)
	_ servertypes.Application = (*SimApp)(nil)
)

// SimApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type SimApp struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry codectypes.InterfaceRegistry

	mm           *module.Manager
	configurator module.Configurator
	keys         map[string]*sdk.KVStoreKey
	mkeys        map[string]*sdk.MemoryStoreKey
	tkeys        map[string]*sdk.TransientStoreKey

	// Cosmos SDK Modules
	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	ParamsKeeper  paramskeeper.Keeper
	StakingKeeper stakingkeeper.Keeper
	UpgradeKeeper upgradekeeper.Keeper
	// IBC Modules
	CapabilityKeeper *capabilitykeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper
	TransferKeeper   transferkeeper.Keeper
	// Custom Modules
	AuraKeeper *aurakeeper.Keeper
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".simapp")
}

// NewSimApp returns a reference to an initialized SimApp.
func NewSimApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *SimApp {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp("SimApp", logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, paramstypes.StoreKey, stakingtypes.StoreKey,
		upgradetypes.ModuleName, capabilitytypes.StoreKey, host.StoreKey, transfertypes.StoreKey,
		auratypes.ModuleName,
	)
	mkeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, keys); err != nil {
		fmt.Printf("failed to load state streaming: %s", err)
		os.Exit(1)
	}

	app := &SimApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
		mkeys:             mkeys,
		tkeys:             tkeys,
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, permissions,
	)

	app.AuraKeeper = aurakeeper.NewKeeper(
		appCodec, keys[auratypes.ModuleName], "ausdy", nil,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.ModuleAccountAddrs(),
	).WithSendCoinsRestriction(app.AuraKeeper.SendRestrictionFn)
	app.AuraKeeper.SetBankKeeper(app.BankKeeper)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(
		make(map[int64]bool), keys[upgradetypes.StoreKey], appCodec, DefaultNodeHome, app.BaseApp,
	)

	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec, keys[capabilitytypes.StoreKey], mkeys[capabilitytypes.MemStoreKey],
	)

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(host.ModuleName)
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[host.StoreKey], app.GetSubspace(host.ModuleName), app.StakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(transfertypes.ModuleName)
	app.TransferKeeper = transferkeeper.NewKeeper(
		appCodec, keys[transfertypes.StoreKey], app.GetSubspace(transfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)

	router := porttypes.NewRouter().AddRoute(transfertypes.ModuleName, transfer.NewIBCModule(app.TransferKeeper))
	app.IBCKeeper.SetRouter(router)

	app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		bank.NewAppModule(appCodec, app.BankKeeper.(bankkeeper.BaseKeeper), app.AccountKeeper),
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		params.NewAppModule(app.ParamsKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		transfer.NewAppModule(app.TransferKeeper),
		aura.NewAppModule(app.AuraKeeper),
	)

	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName, capabilitytypes.ModuleName, stakingtypes.ModuleName,
		host.ModuleName, transfertypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName,
		genutiltypes.ModuleName, paramstypes.ModuleName, auratypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		stakingtypes.ModuleName, host.ModuleName, transfertypes.ModuleName, capabilitytypes.ModuleName,
		authtypes.ModuleName, banktypes.ModuleName, genutiltypes.ModuleName, paramstypes.ModuleName,
		upgradetypes.ModuleName, auratypes.ModuleName,
	)
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName,
		stakingtypes.ModuleName, host.ModuleName, genutiltypes.ModuleName, transfertypes.ModuleName,
		paramstypes.ModuleName, upgradetypes.ModuleName, auratypes.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	app.MountKVStores(keys)
	app.MountMemoryStores(mkeys)
	app.MountTransientStores(tkeys)

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			AuraKeeper:      app.AuraKeeper,
			BankKeeper:      app.BankKeeper,
			IBCKeeper:       app.IBCKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			FeegrantKeeper:  nil,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		panic(err)
	}

	app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

func (app *SimApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *SimApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *SimApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *SimApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *SimApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (app *SimApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range permissions {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *SimApp) SimulationManager() *module.SimulationManager {
	return nil
}

func (app *SimApp) RegisterAPIRoutes(apiSvr *api.Server, _ config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
}

func (app *SimApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *SimApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

func (app *SimApp) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

func (app *SimApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(capabilitytypes.ModuleName)
	paramsKeeper.Subspace(host.ModuleName)
	paramsKeeper.Subspace(transfertypes.ModuleName)

	return paramsKeeper
}
