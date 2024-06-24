# Events

## `burn`

This event is emitted by the bank module whenever coins are burned (including USDY).

```json
{
  "type": "burn",
  "attributes": [
    {
      "key": "burner",
      "value": "noble1d47qadpqs5kjuc5ghc5uglyhtlc4mq0wmc3n33"
    },
    {
      "key": "amount",
      "value": "1000000000000000000ausdy"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgBurn`](./02_messages.md#burn)

## `coinbase`

This event is emitted by the bank module whenever coins are minted (including USDY).

```json
{
  "type": "coinbase",
  "attributes": [
    {
      "key": "minter",
      "value": "noble1d47qadpqs5kjuc5ghc5uglyhtlc4mq0wmc3n33"
    },
    {
      "key": "amount",
      "value": "1000000000000000000ausdy"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgMint`](./02_messages.md#mint)

## Paused

This event is emitted whenever USDY is paused.

```json
{
  "type": "aura.v1.Paused",
  "attributes": [
    {
      "key": "account",
      "value": "noble1signer"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgPause`](./02_messages.md#pause)

## Unpaused

This event is emitted whenever USDY is unpaused.

```json
{
  "type": "aura.v1.Unpaused",
  "attributes": [
    {
      "key": "account",
      "value": "noble1signer"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgUnpause`](./02_messages.md#unpause)

## OwnershipTransferStarted

This event is emitted when an ownership transfer is started.

```json
{
  "type": "aura.v1.OwnershipTransferStarted",
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

- [`aura.v1.MsgTransferOwnership`](./02_messages.md#transfer-ownership)

## OwnershipTransferred

This event is emitted when an ownership transfer is finalized.

```json
{
  "type": "aura.v1.OwnershipTransferred",
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

- [`aura.v1.MsgAcceptOwnership`](./02_messages.md#accept-ownership)

## BurnerAdded

This event is emitted when a new burner is added.

```json
{
  "type": "aura.v1.BurnerAdded",
  "attributes": [
    {
      "key": "address",
      "value": "noble1burner"
    },
    {
      "key": "allowance",
      "value": "0"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgAddBurner`](./02_messages.md#add-burner)

## BurnerRemoved

This event is emitted when a burner is removed.

```json
{
  "type": "aura.v1.BurnerRemoved",
  "attributes": [
    {
      "key": "address",
      "value": "noble1burner"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgRemoveBurner`](./02_messages.md#remove-burner)

## BurnerUpdated

This event is emitted when a burner's allowance is updated.

```json
{
  "type": "aura.v1.BurnerUpdated",
  "attributes": [
    {
      "key": "address",
      "value": "noble1burner"
    },
    {
      "key": "new_allowance",
      "value": "1000000000000000000"
    },
    {
      "key": "previous_allowance",
      "value": "0"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgSetBurnerAllowance`](./02_messages.md#set-burner-allowance)

## MinterAdded

This event is emitted when a new minter is added.

```json
{
  "type": "aura.v1.MinterAdded",
  "attributes": [
    {
      "key": "address",
      "value": "noble1minter"
    },
    {
      "key": "allowance",
      "value": "0"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgAddMinter`](./02_messages.md#add-minter)

## MinterRemoved

This event is emitted when a minter is removed.

```json
{
  "type": "aura.v1.MinterRemoved",
  "attributes": [
    {
      "key": "address",
      "value": "noble1minter"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgRemoveMinter`](./02_messages.md#remove-minter)

## MinterUpdated

This event is emitted when a minter's allowance is updated.

```json
{
  "type": "aura.v1.MinterUpdated",
  "attributes": [
    {
      "key": "address",
      "value": "noble1minter"
    },
    {
      "key": "new_allowance",
      "value": "1000000000000000000"
    },
    {
      "key": "previous_allowance",
      "value": "0"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgSetMinterAllowance`](./02_messages.md#set-minter-allowance)

## PauserAdded

This event is emitted when a new pauser is added.

```json
{
  "type": "aura.v1.PauserAdded",
  "attributes": [
    {
      "key": "address",
      "value": "noble1pauser"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgAddPauser`](./02_messages.md#add-pauser)

## PauserRemoved

This event is emitted when a pauser is removed.

```json
{
  "type": "aura.v1.PauserRemoved",
  "attributes": [
    {
      "key": "address",
      "value": "noble1pauser"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgRemovePauser`](./02_messages.md#remove-pauser)

## ChannelAllowed

This event is emitted when a channel is allowed.

```json
{
  "type": "aura.v1.ChannelAllowed",
  "attributes": [
    {
      "key": "channel",
      "value": "channel-0"
    }
  ]
}
```

This event is emitted by the following transactions:

- [`aura.v1.MsgAllowChannel`](./02_messages.md#allow-channel)
