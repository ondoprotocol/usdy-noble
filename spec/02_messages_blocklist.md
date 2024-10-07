# Blocklist Messages

## Transfer Ownership

`aura.blocklist.v1.MsgTransferOwnership`

A message that initiates an ownership transfer to a provided address.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.blocklist.v1.MsgTransferOwnership",
        "signer": "noble1signer",
        "new_owner": "noble1owner"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `new_owner` — The Noble address to initiate an ownership transfer to.

### Requirements

- Signer must be the current [`owner`](./01_state_blocklist.md#owner).
- `new_owner` must not be the current [`owner`](./01_state_blocklist.md#owner).

### State Changes

- [`pending_owner`](./01_state_blocklist.md#pending-owner)

### Events Emitted

- [`aura.blocklist.v1.OwnershipTransferStarted`](./03_events_blocklist#ownershiptransferstarted)

## Accept Ownership

`aura.blocklist.v1.MsgAcceptOwnership`

A message that finalizes an ownership transfer.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.blocklist.v1.MsgAcceptOwnership",
        "signer": "noble1owner"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
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

- [`aura.blocklist.v1.OwnershipTransferred`](./03_events_blocklist.md#ownershiptransferred)

## Add To Blocklist

`aura.blocklist.v1.MsgAddToBlocklist`

A message that adds addresses to Aura's blocklist.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.blocklist.v1.MsgAddToBlocklist",
        "signer": "noble1signer",
        "accounts": [
          "noble1alice",
          "noble1bob",
          "noble1charlie"
        ]
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `accounts` — A list of Noble address to add to the blocklist.

### Requirements

- Signer must be the current [`owner`](./01_state_blocklist.md#owner).

### State Changes

- [`blocked_address`](./01_state_blocklist.md#blocked-addresses)

### Events Emitted

- [`aura.blocklist.v1.BlockedAddressesAdded`](./03_events_blocklist#blockedaddressesadded)

## Remove From Blocklist

`aura.blocklist.v1.MsgRemoveFromBlocklist`

A message that removes addresses from Aura's blocklist.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.blocklist.v1.MsgRemoveFromBlocklist",
        "signer": "noble1signer",
        "accounts": [
          "noble1alice",
          "noble1bob",
          "noble1charlie"
        ]
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

### Arguments

- `accounts` — A list of Noble address to remove from the blocklist.

### Requirements

- Signer must be the current [`owner`](./01_state_blocklist.md#owner).

### State Changes

- [`blocked_address`](./01_state_blocklist.md#blocked-addresses)

### Events Emitted

- [`aura.blocklist.v1.BlockedAddressesRemoved`](./03_events_blocklist#blockedaddressesremoved)
