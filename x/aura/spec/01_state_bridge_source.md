# Bridge Source State

## Paused

The paused field is of type boolean.
It is used to store the current paused state of the USDY source bridge.

```go
var PausedKey = []byte("bridge/source/paused")
```

It is updated by the following messages:

- `aura.bridge.source.v1.MsgPause`
- `aura.bridge.source.v1.MsgUnpause`

## Owner

The owner field is of type string, specifically a Noble encoded address.
It is used to store the current owner of the blocklist submodule.

```go
var OwnerKey = []byte("bridge/source/owner")
```

It is updated by the following messages:

- `aura.bridge.source.v1.MsgTransferOwnership`

## Nonce

The nonce field is of type uint64.
It is used to store the next nonce to be used in a source transfer.

```go
var NonceKey = []byte("bridge/source/nonce")
```

It is updated by the following messages:

- `aura.bridge.source.v1.MsgBurnAndCallAxelar`

## Channel

The channel field is of type string.
It is used to store the enshrined IBC connection to Axelar.

```go
var ChannelKey = []byte("bridge/source/channel")
```

It is updated by the following message:

- TODO

## Destinations

The destinations field is a mapping between string and string.
It is used to store the destination contracts of USDY transfers on other chains.

```go
var MinterPrefix = []byte("minter/")
```

It is updated by the following messages:

- TODO
