# Blocklist Events

## OwnershipTransferStarted

This event is emitted when an ownership transfer is started.

```json
{
  "type": "aura.blocklist.v1.OwnershipTransferStarted",
  "attributes": [
    {
      "key": "new_owner",
      "value": "noble1owner"
    },
    {
      "key": "previous_owner",
      "value": "noble1signer"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.blocklist.v1.MsgTransferOwnership`](./02_messages_blocklist.md#transfer-ownership)

## OwnershipTransferred

This event is emitted when an ownership transfer is finalized.

```json
{
  "type": "aura.blocklist.v1.OwnershipTransferred",
  "attributes": [
    {
      "key": "new_owner",
      "value": "noble1owner"
    },
    {
      "key": "previous_owner",
      "value": "noble1signer"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.blocklist.v1.MsgAcceptOwnership`](./02_messages_blocklist.md#accept-ownership)

## BlockedAddressesAdded

This event is emitted whenever addresses are added to Aura's blocklist.

```json
{
  "type": "aura.blocklist.v1.BlockedAddressesAdded",
  "attributes": [
    {
      "key": "accounts",
      "value": ["noble1alice","noble1bob","noble1charlie"]
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.blocklist.v1.MsgAddToBlocklist`](./02_messages_blocklist.md#add-to-blocklist)

## BlockedAddressesRemoved

This event is emitted whenever addresses are removed from Aura's blocklist.

```json
{
  "type": "aura.blocklist.v1.BlockedAddressesRemoved",
  "attributes": [
    {
      "key": "accounts",
      "value": ["noble1alice","noble1bob","noble1alice"]
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.blocklist.v1.MsgRemoveFromBlocklist`](./02_messages_blocklist.md#remove-from-blocklist)
