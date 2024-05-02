# Blocklist State

## Owner

The owner field is of type string, specifically a Noble encoded address, stored
via an [`collections.Item`][item]. It is used to store the current owner of the
blocklist submodule.

```go
var OwnerKey = []byte("blocklist/owner")
```

It is updated by the following messages:

- [`aura.blocklist.v1.MsgAcceptOwnership`](./02_messages_blocklist.md#accept-ownership)

## Pending Owner

The pending owner field is of type string, specifically a Noble encoded address,
stored via an [`collections.Item`][item]. It is used to store the current
pending owner of the blocklist submodule.

```go
var PendingOwnerKey = []byte("blocklist/pending_owner")
```

It is updated by the following messages:

- [`aura.blocklist.v1.MsgTransferOwnership`](./02_messages_blocklist.md#transfer-ownership)
- [`aura.blocklist.v1.MsgAcceptOwnership`](./02_messages_blocklist.md#accept-ownership)

## Blocked Addresses

The blocked addresses field is a mapping between string (a Noble encoded
address) and boolean, stored via a [`collections.Map`][map]. It is used to store
all blocked addresses that can't interact with USDY.

```go
var BlockedAddressPrefix = []byte("blocklist/blocked_address/")
```

It is updated by the following messages:

- [`aura.blocklist.v1.MsgAddToBlocklist`](./02_messages_blocklist.md#add-to-blocklist)
- [`aura.blocklist.v1.MsgRemoveFromBlocklist`](./02_messages_blocklist.md#remove-from-blocklist)

[item]: https://docs.cosmos.network/main/build/packages/collections#item

[map]: https://docs.cosmos.network/main/build/packages/collections#map

[set]: https://docs.cosmos.network/main/build/packages/collections#keyset
