# Messages

## Burn

`aura.v1.MsgBurn`

A message that burns USDY from a user.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgBurn",
        "signer": "noble1signer",
        "from": "noble1user",
        "amount": "1000000000000000000"
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

- `from` — The Noble address to burn USDY from.
- `amount` — The amount of USDY to burn.

### Requirements

- Signer must be one of the allowed [`burners`](./01_state.md#burners).
- Burner must have enough allowance.

### State Changes

- [`burners`](./01_state.md#burners)
- User balance via `x/bank` module.

### Events Emitted

- [`burn`](./03_events.md#burn)

## Mint

`aura.v1.MsgMint`

A message that mints USDY to a user.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgMint",
        "signer": "noble1signer",
        "to": "noble1user",
        "amount": "1000000000000000000"
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

- `to` — The Noble address to mint USDY to.
- `amount` — The amount of USDY to mint.

### Requirements

- Signer must be one of the allowed [`minters`](./01_state.md#minters).
- Minter must have enough allowance.

### State Changes

- [`minters`](./01_state.md#minters)
- User balance via `x/bank` module.

### Events Emitted

- [`coinbase`](./03_events.md#coinbase)

## Pause

`aura.v1.MsgPause`

A message that pauses USDY.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgPause",
        "signer": "noble1signer"
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

- Signer must be one of the allowed [`pausers`](./01_state.md#pausers).

### State Changes

- [`paused`](./01_state.md#paused)

### Events Emitted

- [`aura.v1.Paused`](./03_events.md#paused)

## Unpause

`aura.v1.MsgUnpause`

A message that unpauses USDY.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgUnpause",
        "signer": "noble1signer"
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

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`paused`](./01_state.md#paused)

### Events Emitted

- [`aura.v1.Unpaused`](./03_events.md#unpaused)

## Transfer Ownership

`aura.v1.MsgTransferOwnership`

A message that initiates an ownership transfer to a provided address.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgTransferOwnership",
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

- Signer must be the current [`owner`](./01_state.md#owner).
- `new_owner` must not be the current [`owner`](./01_state.md#owner).

### State Changes

- [`pending_owner`](./01_state.md#pending-owner)

### Events Emitted

- [`OwnershipTransferStarted`](./03_events.md#ownershiptransferstarted)

## Accept Ownership

`aura.v1.MsgAcceptOwnership`

A message that finalizes an ownership transfer.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgAcceptOwnership",
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

- [`pending_owner`](./01_state.md#pending-owner) must be set in state first,
  initiated via
  a [`MsgTransferOwnership`](#transfer-ownership) message being previously
  executed.
- Signer must be the current [`pending_owner`](./01_state.md#pending-owner).

### State Changes

- [`owner`](./01_state.md#owner)
- [`pending_owner`](./01_state.md#pending-owner)

### Events Emitted

- [`aura.v1.OwnershipTransferred`](./03_events.md#ownershiptransferred)

## Add Burner

`aura.v1.MsgAddBurner`

A message that adds a new burner, with an initial allowance.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgAddBurner",
        "signer": "noble1signer",
        "burner": "noble1burner",
        "allowance": "0"
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

- `burner` — The Noble address to add as a new burner.
- `allowance` — The initial burn allowance of this new burner.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`burners`](./01_state.md#burners)

### Events Emitted

- [`aura.v1.BurnerAdded`](./03_events.md#burneradded)

## Remove Burner

`aura.v1.MsgRemoveBurner`

A message that removes a burner.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgRemoveBurner",
        "signer": "noble1signer",
        "burner": "noble1burner"
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

- `burner` — The Noble address to remove as a burner.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`burners`](./01_state.md#burners)

### Events Emitted

- [`aura.v1.BurnerRemoved`](./03_events.md#burnerremoved)

## Set Burner Allowance

`aura.v1.MsgSetBurnerAllowance`

A message that sets the allowance of a burner.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgSetBurnerAllowance",
        "signer": "noble1signer",
        "burner": "noble1burner",
        "allowance": "1000000000000000000"
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

- `burner` — The Noble address to update the burn allowance for.
- `allowance` — The burn allowance to update to.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`burners`](./01_state.md#burners)

### Events Emitted

- [`aura.v1.BurnerUpdated`](./03_events.md#burnerupdated)

## Add Minter

`aura.v1.MsgAddMinter`

A message that adds a new minter, with an initial allowance.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgAddMinter",
        "signer": "noble1signer",
        "minter": "noble1minter",
        "allowance": "0"
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

- `minter` — The Noble address to add as a new minter.
- `allowance` — The initial mint allowance of this new minter.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`minters`](./01_state.md#minters)

### Events Emitted

- [`aura.v1.MinterAdded`](./03_events.md#minteradded)

## Remove Minter

`aura.v1.MsgRemoveMinter`

A message that removes a minter.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgRemoveMinter",
        "signer": "noble1signer",
        "minter": "noble1minter"
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

- `minter` — The Noble address to remove as a minter.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`minters`](./01_state.md#minters)

### Events Emitted

- [`aura.v1.MinterRemoved`](./03_events.md#minterremoved)

## Set Minter Allowance

`aura.v1.MsgSetMinterAllowance`

A message that sets the allowance of a minter.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgSetMinterAllowance",
        "signer": "noble1signer",
        "minter": "noble1minter",
        "allowance": "1000000000000000000"
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

- `minter` — The Noble address to update the mint allowance for.
- `allowance` — The mint allowance to update to.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`minters`](./01_state.md#minters)

### Events Emitted

- [`aura.v1.MinterUpdated`](./03_events.md#minterupdated)

## Add Pauser

`aura.v1.MsgAddPauser`

A message that adds a new pauser.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgAddPauser",
        "signer": "noble1signer",
        "pauser": "noble1pauser"
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

- `pauser` — The Noble address to add as a new pauser.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`pausers`](./01_state.md#pausers)

### Events Emitted

- [`aura.v1.PauserAdded`](./03_events.md#pauseradded)

## Remove Pauser

`aura.v1.MsgRemovePauser`

A message that removes a pauser.

```json
{
  "body": {
    "messages": [
      {
        "@type": "/aura.v1.MsgRemovePauser",
        "signer": "noble1signer",
        "pauser": "noble1pauser"
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

- `pauser` — The Noble address to remove as a pauser.

### Requirements

- Signer must be the current [`owner`](./01_state.md#owner).

### State Changes

- [`pausers`](./01_state.md#pausers)

### Events Emitted

- [`aura.v1.PauserRemoved`](./03_events.md#pauserremoved)
