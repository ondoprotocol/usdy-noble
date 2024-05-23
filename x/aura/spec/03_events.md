# Events

## Paused

This event is emitted whenever USDY is paused.

```shell
- attributes:
  - index: true
    key: account
    value: '"noble1signer"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.Paused
```

This event is emitted by the following transactions:

- [`aura.v1.MsgPause`](./02_messages.md#pause)

## Unpaused

This event is emitted whenever USDY is unpaused.

```shell
- attributes:
  - index: true
    key: account
    value: '"noble1signer"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.Unpaused
```

This event is emitted by the following transactions:

- [`aura.v1.MsgUnpause`](./02_messages.md#unpause)

## OwnershipTransferStarted

This event is emitted when an ownership transfer is started.

```shell
- attributes:
  - index: true
    key: new_owner
    value: '"noble1owner"'
  - index: true
    key: previous_owner
    value: '"noble1signer"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.OwnershipTransferStarted
```

This event is emitted by the following transactions:

- [`aura.v1.MsgTransferOwnership`](./02_messages.md#transfer-ownership)

## OwnershipTransferred

This event is emitted when an ownership transfer is finalized.

```shell
- attributes:
  - index: true
    key: new_owner
    value: '"noble1owner"'
  - index: true
    key: previous_owner
    value: '"noble1signer"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.OwnershipTransferred
```

This event is emitted by the following transactions:

- [`aura.v1.MsgAcceptOwnership`](./02_messages.md#accept-ownership)

## BurnerAdded

This event is emitted when a new burner is added.

```shell
- attributes:
  - index: true
    key: address
    value: '"noble1burner"'
  - index: true
    key: allowance
    value: '"0"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.BurnerAdded
```

This event is emitted by the following transactions:

- [`aura.v1.MsgAddBurner`](./02_messages.md#add-burner)

## BurnerRemoved

This event is emitted when a burner is removed.

```shell
- attributes:
  - index: true
    key: address
    value: '"noble1burner"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.BurnerRemoved
```

This event is emitted by the following transactions:

- [`aura.v1.MsgRemoveBurner`](./02_messages.md#remove-burner)

## BurnerUpdated

This event is emitted when a burner's allowance is updated.

```shell
- attributes:
  - index: true
    key: address
    value: '"noble1burner"'
  - index: true
    key: new_allowance
    value: '"1000000000000000000"'
  - index: true
    key: previous_allowance
    value: '"0"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.BurnerUpdated
```

This event is emitted by the following transactions:

- [`aura.v1.MsgSetBurnerAllowance`](./02_messages.md#set-burner-allowance)

## MinterAdded

This event is emitted when a new minter is added.

```shell
- attributes:
  - index: true
    key: address
    value: '"noble1minter"'
  - index: true
    key: allowance
    value: '"0"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.MinterAdded
```

This event is emitted by the following transactions:

- [`aura.v1.MsgAddMinter`](./02_messages.md#add-minter)

## MinterRemoved

This event is emitted when a minter is removed.

```shell
- attributes:
  - index: true
    key: address
    value: '"noble1minter"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.MinterRemoved
```

This event is emitted by the following transactions:

- [`aura.v1.MsgRemoveMinter`](./02_messages.md#remove-minter)

## MinterUpdated

This event is emitted when a minter's allowance is updated.

```shell
- attributes:
  - index: true
    key: address
    value: '"noble1minter"'
  - index: true
    key: new_allowance
    value: '"1000000000000000000"'
  - index: true
    key: previous_allowance
    value: '"0"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.MinterUpdated
```

This event is emitted by the following transactions:

- [`aura.v1.MsgSetMinterAllowance`](./02_messages.md#set-minter-allowance)

## PauserAdded

This event is emitted when a new pauser is added.

```shell
- attributes:
  - index: true
    key: address
    value: '"noble1pauser"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.PauserAdded
```

This event is emitted by the following transactions:

- [`aura.v1.MsgAddPauser`](./02_messages.md#add-pauser)

## PauserRemoved

This event is emitted when a pauser is removed.

```shell
- attributes:
  - index: true
    key: address
    value: '"noble1pauser"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.PauserRemoved
```

This event is emitted by the following transactions:

- [`aura.v1.MsgRemovePauser`](./02_messages.md#remove-pauser)
