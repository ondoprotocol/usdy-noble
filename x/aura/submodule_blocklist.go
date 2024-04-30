package aura

import (
	"context"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	modulev1 "github.com/noble-assets/aura/api/aura/blocklist/module/v1"
	blocklistv1 "github.com/noble-assets/aura/api/aura/blocklist/v1"
	"github.com/noble-assets/aura/x/aura/keeper"
	"github.com/noble-assets/aura/x/aura/types/blocklist"
)

var (
	_ module.AppModuleBasic = BlocklistSubmodule{}
	_ appmodule.AppModule   = BlocklistSubmodule{}
	_ module.HasServices    = BlocklistSubmodule{}
)

//

type BlocklistSubmoduleBasic struct{}

func NewBlocklistSubmoduleBasic() BlocklistSubmoduleBasic {
	return BlocklistSubmoduleBasic{}
}

func (BlocklistSubmoduleBasic) Name() string { return blocklist.SubmoduleName }

func (BlocklistSubmoduleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	blocklist.RegisterLegacyAminoCodec(cdc)
}

func (BlocklistSubmoduleBasic) RegisterInterfaces(reg codectypes.InterfaceRegistry) {
	blocklist.RegisterInterfaces(reg)
}

func (BlocklistSubmoduleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := blocklist.RegisterQueryHandlerClient(context.Background(), mux, blocklist.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

//

type BlocklistSubmodule struct {
	BlocklistSubmoduleBasic

	keeper *keeper.Keeper
}

func NewBlocklistSubmodule(keeper *keeper.Keeper) BlocklistSubmodule {
	return BlocklistSubmodule{
		BlocklistSubmoduleBasic: NewBlocklistSubmoduleBasic(),
		keeper:                  keeper,
	}
}

func (BlocklistSubmodule) IsOnePerModuleType() {}

func (BlocklistSubmodule) IsAppModule() {}

func (m BlocklistSubmodule) RegisterServices(cfg module.Configurator) {
	blocklist.RegisterMsgServer(cfg.MsgServer(), keeper.NewBlocklistMsgServer(m.keeper))
	blocklist.RegisterQueryServer(cfg.QueryServer(), keeper.NewBlocklistQueryServer(m.keeper))
}

//

func (BlocklistSubmodule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: blocklistv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod:      "TransferOwnership",
					Use:            "transfer-ownership [new-owner]",
					Short:          "Transfer ownership of submodule",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "new_owner"}},
				},
				{
					RpcMethod: "AcceptOwnership",
					Use:       "accept-ownership",
					Short:     "Accept ownership of submodule",
					Long:      "Accept ownership of submodule, assuming there is an pending ownership transfer",
				},
				{
					RpcMethod: "AddToBlocklist",
					Use:       "add-to-blocklist [addresses ...]",
					Short:     "Add a list of accounts to the blocklist",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{
						ProtoField: "accounts",
						Varargs:    true,
					}},
				},
				{
					RpcMethod: "RemoveFromBlocklist",
					Use:       "remove-from-blocklist [addresses ...]",
					Short:     "Remove a list of accounts from the blocklist",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{
						ProtoField: "accounts",
						Varargs:    true,
					}},
				},
			},
		},
		Query: &autocliv1.ServiceCommandDescriptor{
			Service:           blocklistv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				// NOTE(@john): Queries are simple, so no need to configure.
			},
		},
	}
}

//

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideBlocklistModule),
	)
}

type BlocklistModuleInputs struct {
	depinject.In

	Config *modulev1.Module
	Keeper *keeper.Keeper
}

type BlocklistModuleOutputs struct {
	depinject.Out

	Module appmodule.AppModule
}

func ProvideBlocklistModule(in BlocklistModuleInputs) BlocklistModuleOutputs {
	m := NewBlocklistSubmodule(in.Keeper)

	return BlocklistModuleOutputs{Module: m}
}
