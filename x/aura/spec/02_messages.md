# Messages

## Burn

`aura.v1.MsgBurn`

A message that burns USDY from a user.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgBurn
    amount: "1000000000000000000"
    from: noble1user
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `from` — The Noble address to burn USDY from.
- `amount` — The amount of USDY to burn.

### Requirements

- Signer must be one of the allowed [`burners`](./01_state.md#burners).
- Burner must have enough allowance.

### State Changes

- [`burners`](./01_state.md#burners)
- User balance via `x/bank` module.

### Events Emitted

This message emits no events.

## Mint

`aura.v1.MsgMint`

A message that mints USDY to a user.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgMint
    amount: "1000000000000000000"
    signer: noble1signer
    to: noble1user
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `to` — The Noble address to mint USDY to.
- `amount` — The amount of USDY to mint.

### Requirements

- Signer must be one of the allowed [`minters`](./01_state.md#minters).
- Minter must have enough allowance.

### State Changes

- [`minters`](./01_state.md#minters)
- User balance via `x/bank` module.

### Events Emitted

This message emits no events.

## Pause

`aura.v1.MsgPause`

A message that pauses USDY.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgPause
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

This message takes no arguments.

### Requirements

- Signer must be one of the allowed [`pausers`](./01_state.md#pausers).

### State Changes

- [`paused`](./01_state.md#paused)

### Events Emitted

- [`aura.v1.Paused`](./03_events.md#paused)

## Unpause

`aura.v1.MsgUnpause`

A message that unpauses USDY.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgUnpause
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

This message takes no arguments.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`paused`](./01_state.md#paused)

### Events Emitted

- [`aura.v1.Unpaused`](./03_events.md#unpaused)

## Transfer Ownership

`aura.v1.MsgTransferOwnership`

A message that initiates an ownership transfer to a provided address.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgTransferOwnership
    new_owner: noble1owner
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `new_owner` — The Noble address to initiate an ownership transfer to.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`pending_owner`](./01_state.md#pending-owner)

### Events Emitted

- [`OwnershipTransferStarted`](./03_events.md#ownershiptransferstarted)

## Accept Ownership

`aura.v1.MsgAcceptOwnership`

A message that finalizes an ownership transfer.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgAcceptOwnership
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

This message takes no arguments.

### Requirements

- [`pending_owner`](./01_state.md#pending-owner) must be set in state first,
  initiated via
  a [`MsgTransferOwnership`](#transfer-ownership) message being previously
  executed.
- Signer must be the current [`pending_owner`](./01_state.md#pending-owner).

### State Changes

- [`owner`](./01_state.md#owner)
- [`pending_owner`](./01_state.md#pending-owner)

### Events Emitted

- [`aura.v1.OwnershipTransferred`](./03_events.md#ownershiptransferred)

## Add Burner

`aura.v1.MsgAddBurner`

A message that adds a new burner, with an initial allowance.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgAddBurner
    allowance: "0"
    burner: noble1burner
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `burner` — The Noble address to add as a new burner.
- `allowance` — The initial burn allowance of this new burner.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`burners`](./01_state.md#burners)

### Events Emitted

- [`aura.v1.BurnerAdded`](./03_events.md#burneradded)

## Remove Burner

`aura.v1.MsgRemoveBurner`

A message that removes a burner.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgRemoveBurner
    burner: noble1burner
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `burner` — The Noble address to remove as a burner.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`burners`](./01_state.md#burners)

### Events Emitted

- [`aura.v1.BurnerRemoved`](./03_events.md#burnerremoved)

## Set Burner Allowance

`aura.v1.MsgSetBurnerAllowance`

A message that sets the allowance of a burner.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgSetBurnerAllowance
    allowance: "1000000000000000000"
    burner: noble1burner
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `burner` — The Noble address to update the burn allowance for.
- `allowance` — The burn allowance to update to.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`burners`](./01_state.md#burners)

### Events Emitted

- [`aura.v1.BurnerUpdated`](./03_events.md#burnerupdated)

## Add Minter

`aura.v1.MsgAddMinter`

A message that adds a new minter, with an initial allowance.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgAddMinter
    allowance: "0"
    minter: noble1minter
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `minter` — The Noble address to add as a new minter.
- `allowance` — The initial mint allowance of this new minter.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`minters`](./01_state.md#minters)

### Events Emitted

- [`aura.v1.MinterAdded`](./03_events.md#minteradded)

## Remove Minter

`aura.v1.MsgRemoveMinter`

A message that removes a minter.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgRemoveMinter
    minter: noble1minter
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `minter` — The Noble address to remove as a minter.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`minters`](./01_state.md#minters)

### Events Emitted

- [`aura.v1.MinterRemoved`](./03_events.md#minterremoved)

## Set Minter Allowance

`aura.v1.MsgSetMinterAllowance`

A message that sets the allowance of a minter.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgSetMinterAllowance
    allowance: "1000000000000000000"
    minter: noble1minter
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `minter` — The Noble address to update the mint allowance for.
- `allowance` — The mint allowance to update to.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`minters`](./01_state.md#minters)

### Events Emitted

- [`aura.v1.MinterUpdated`](./03_events.md#minterupdated)

## Add Pauser

`aura.v1.MsgAddPauser`

A message that adds a new pauser.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgAddPauser
    pauser: noble1pauser
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `pauser` — The Noble address to add as a new pauser.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`pausers`](./01_state.md#pausers)

### Events Emitted

- [`aura.v1.PauserAdded`](./03_events.md#pauseradded)

## Remove Pauser

`aura.v1.MsgRemovePauser`

A message that removes a pauser.

```shell
auth_info:
  fee:
    amount: []
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /aura.v1.MsgRemovePauser
    pauser: noble1pauser
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `pauser` — The Noble address to remove as a pauser.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`pausers`](./01_state.md#pausers)

### Events Emitted

- [`aura.v1.PauserRemoved`](./03_events.md#pauserremoved)
