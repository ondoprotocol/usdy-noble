package keeper_test

import (
	"testing"

	"github.com/noble-assets/ondo/utils"
	"github.com/noble-assets/ondo/utils/mocks"
	"github.com/noble-assets/ondo/x/usdy/keeper"
	"github.com/noble-assets/ondo/x/usdy/types"
	"github.com/stretchr/testify/require"
)

func TestPause(t *testing.T) {
	k, ctx := mocks.USDYKeeper(t)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to pause with no pauser set.
	_, err := server.Pause(ctx, &types.MsgPause{})
	// ASSERT: The action should've failed due to no pauser set.
	require.ErrorContains(t, err, "unable to retrieve pauser from state")
	paused, _ := k.Paused.Get(ctx)
	require.False(t, paused)

	// ARRANGE: Set pauser in state.
	pauser := utils.TestAccount()
	require.NoError(t, k.Pauser.Set(ctx, pauser.Address))

	// ACT: Attempt to pause with invalid signer.
	_, err = server.Pause(ctx, &types.MsgPause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPauser.Error())
	paused, _ = k.Paused.Get(ctx)
	require.False(t, paused)

	// ACT: Attempt to pause.
	_, err = server.Pause(ctx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	paused, _ = k.Paused.Get(ctx)
	require.True(t, paused)

	// ACT: Attempt to pause again.
	_, err = server.Pause(ctx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've failed due to module being paused already.
	require.ErrorContains(t, err, "module is already paused")
	paused, _ = k.Paused.Get(ctx)
	require.True(t, paused)
}

func TestUnpause(t *testing.T) {
	k, ctx := mocks.USDYKeeper(t)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set paused state to true.
	require.NoError(t, k.Paused.Set(ctx, true))

	// ACT: Attempt to unpause with no pauser set.
	_, err := server.Unpause(ctx, &types.MsgUnpause{})
	// ASSERT: The action should've failed due to no pauser set.
	require.ErrorContains(t, err, "unable to retrieve pauser from state")
	paused, _ := k.Paused.Get(ctx)
	require.True(t, paused)

	// ARRANGE: Set pauser in state.
	pauser := utils.TestAccount()
	require.NoError(t, k.Pauser.Set(ctx, pauser.Address))

	// ACT: Attempt to unpause with invalid signer.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPauser.Error())
	paused, _ = k.Paused.Get(ctx)
	require.True(t, paused)

	// ACT: Attempt to unpause.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	paused, _ = k.Paused.Get(ctx)
	require.False(t, paused)

	// ACT: Attempt to unpause again.
	_, err = server.Unpause(ctx, &types.MsgUnpause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've failed due to module being unpaused already.
	require.ErrorContains(t, err, "module is already unpaused")
	paused, _ = k.Paused.Get(ctx)
	require.False(t, paused)
}
