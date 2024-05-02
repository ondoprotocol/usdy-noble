# Blocklist Messages

## Transfer Ownership

`aura.blocklist.v1.MsgTransferOwnership`

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
  - '@type': /aura.blocklist.v1.MsgTransferOwnership
    new_owner: noble1owner
    signer: noble1signer
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `new_owner` — The Noble address to initiate an ownership transfer to.

### Requirements

- Signer must be the current [`owner`](./01_state_blocklist.md#owner).

### State Changes

- [`pending_owner`](./01_state_blocklist.md#pending-owner)

### Events Emitted

- [`OwnershipTransferStarted`](./03_events_blocklist#ownershiptransferstarted)

## Accept Ownership

`aura.blocklist.v1.MsgAcceptOwnership`

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
  - '@type': /aura.blocklist.v1.MsgAcceptOwnership
    signer: noble1demo # Noble address of new owner.
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

This message takes no arguments.

### Requirements

- [`pending_owner`](./01_state_blocklist.md#pending-owner) must be set in state
  first, initiated via a [`MsgTransferOwnership`](#transfer-ownership) message
  being previously executed.
- Signer must be the
  current [`pending_owner`](./01_state_blocklist.md#pending-owner).

### State Changes

- [`owner`](./01_state_blocklist.md#owner)
- [`pending_owner`](./01_state_blocklist.md#pending-owner)

### Events Emitted

This message emits no events.

## Add To Blocklist

`aura.blocklist.v1.MsgAddToBlocklist`

A message that adds addresses to Aura's blocklist.

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
  - '@type': /aura.blocklist.v1.MsgAddToBlocklist
    accounts:
    - noble1alice
    - noble1bob
    - noble1charlie
    signer: noble1demo
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `accounts` — A list of Noble address to add to the blocklist.

### Requirements

- Signer must be the current [`owner`](./01_state_blocklist.md#owner).

### State Changes

- [`blocked_address`](./01_state_blocklist.md#blocked-addresses)

### Events Emitted

- [`BlockedAddressesAdded`](./03_events_blocklist#blockedaddressesadded)

## Remove From Blocklist

`aura.blocklist.v1.MsgRemoveFromBlocklist`

A message that removes addresses from Aura's blocklist.

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
  - '@type': /aura.blocklist.v1.MsgRemoveFromBlocklist
    accounts:
    - noble1alice
    - noble1bob
    - noble1charlie
    signer: noble1demo
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
```

### Arguments

- `accounts` — A list of Noble address to remove from the blocklist.

### Requirements

- Signer must be the current [`owner`](./01_state_blocklist.md#owner).

### State Changes

- [`blocked_address`](./01_state_blocklist.md#blocked-addresses)

### Events Emitted

- [`BlockedAddressesRemoved`](./03_events_blocklist#blockedaddressesremoved)
