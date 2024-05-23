# Blocklist Events

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
  type: aura.blocklist.v1.OwnershipTransferStarted
```

This event is emitted by the following transactions:

- [`aura.blocklist.v1.MsgTransferOwnership`](./02_messages_blocklist.md#transfer-ownership)

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
  type: aura.blocklist.v1.OwnershipTransferred
```

This event is emitted by the following transactions:

- [`aura.blocklist.v1.MsgAcceptOwnership`](./02_messages_blocklist.md#accept-ownership)

## BlockedAddressesAdded

This event is emitted whenever addresses are added to Aura's blocklist.

```shell
- attributes:
  - index: true
    key: accounts
    value: '["noble1alice","noble1bob","noble1charlie"]'
  - index: true
    key: msg_index
    value: "0"
  type: aura.blocklist.v1.BlockedAddressesAdded
```

This event is emitted by the following transactions:

- [`aura.blocklist.v1.MsgAddToBlocklist`](./02_messages_blocklist.md#add-to-blocklist)

## BlockedAddressesRemoved

This event is emitted whenever addresses are removed from Aura's blocklist.

```shell
- attributes:
  - index: true
    key: accounts
    value: '["noble1alice","noble1bob","noble1charlie"]'
  - index: true
    key: msg_index
    value: "0"
  type: aura.blocklist.v1.BlockedAddressesRemoved
```

This event is emitted by the following transactions:

- [`aura.blocklist.v1.MsgRemoveFromBlocklist`](./02_messages_blocklist.md#remove-from-blocklist)
