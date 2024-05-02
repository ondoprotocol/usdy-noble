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

This event is emitted when an ownership transfer process is started.

```shell
- attributes:
  - index: true
    key: new_owner
    value: '"noble1owner"'
  - index: true
    key: old_owner
    value: '"noble1signer"'
  - index: true
    key: msg_index
    value: "0"
  type: aura.v1.OwnershipTransferStarted
```

This event is emitted by the following transactions:

- [`aura.v1.MsgTransferOwnership`](./02_messages.md#transfer-ownership)
