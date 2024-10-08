# State

## Paused

The paused field is of type boolean.
It is used to store the current paused state of USDY.

```go
var PausedKey = []byte("paused")
```

It is updated by the following messages:

- [`aura.v1.MsgPause`](./02_messages.md#pause)
- [`aura.v1.MsgUnpause`](./02_messages.md#unpause)

## Owner

The owner field is of type string, specifically a Noble encoded address.
It is used to store the current owner of the module.

```go
var OwnerKey = []byte("owner")
```

It is updated by the following messages:

- [`aura.v1.MsgAcceptOwnership`](./02_messages.md#accept-ownership)

## Pending Owner

The pending owner field is of type string, specifically a Noble encoded address.
It is used to store the current pending owner of the module.

```go
var PendingOwnerKey = []byte("pending_owner")
```

It is updated by the following messages:

- [`aura.v1.MsgTransferOwnership`](./02_messages.md#transfer-ownership)
- [`aura.v1.MsgAcceptOwnership`](./02_messages.md#accept-ownership)

## Burners

The burners field is a mapping between string (a Noble encoded address) and `math.Int`.
It is used to store all burners of USDY, and their current burn allowance.

```go
var BurnerPrefix = []byte("burner/")
```

It is updated by the following messages:

- [`aura.v1.MsgBurn`](./02_messages.md#burn)
- [`aura.v1.MsgAddBurner`](./02_messages.md#add-burner)
- [`aura.v1.MsgRemoveBurner`](./02_messages.md#remove-burner)
- [`aura.v1.MsgSetBurnAllowance`](./02_messages.md#set-burner-allowance)

## Minters

The minters field is a mapping between string (a Noble encoded address) and `math.Int`.
It is used to store all minters of USDY, and their current mint allowance.

```go
var MinterPrefix = []byte("minter/")
```

It is updated by the following messages:

- [`aura.v1.MsgMint`](./02_messages.md#mint)
- [`aura.v1.MsgAddMinter`](./02_messages.md#add-minter)
- [`aura.v1.MsgRemoveMinter`](./02_messages.md#remove-minter)
- [`aura.v1.MsgSetMintAllowance`](./02_messages.md#set-minter-allowance)

## Pausers

The pausers field is a unique set of strings, specifically Noble encoded addresses.
It is used to store all pausers of USDY.

```go
var PauserPrefix = []byte("pauser/")
```

It is updated by the following messages:

- [`aura.v1.MsgAddPauser`](./02_messages.md#add-pauser)
- [`aura.v1.MsgRemovePauser`](./02_messages.md#remove-pauser)

## Blocked Channels

The blocked channels field is a unique set of strings, specifically IBC channels.
It is used to store all blocked IBC transfer channels for USDY.

```go
var BlockedChannelPrefix = []byte("blocked_channel/")
```

It is updated by the following messages:

- [`aura.v1.MsgAddBlockedChannel`](./02_messages.md#add-blocked-channel)
- [`aura.v1.MsgRemoveBlockedChannel`](./02_messages.md#remove-blocked-channel)
