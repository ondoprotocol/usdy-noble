alias ondod=./simapp/build/simd

for arg in "$@"
do
    case $arg in
        -r|--reset)
        rm -rf .ondo
        shift
        ;;
    esac
done

if ! [ -f .ondo/data/priv_validator_state.json ]; then
  ondod init validator --chain-id "ondo-1" --home .ondo &> /dev/null

  ondod keys add validator --home .ondo --keyring-backend test &> /dev/null
  ondod genesis add-genesis-account validator 1000000ustake --home .ondo --keyring-backend test
  BURNER=$(ondod keys add burner --home .ondo --keyring-backend test --output json | jq .address)
  ondod genesis add-genesis-account burner 10000000uusdc --home .ondo --keyring-backend test
  MINTER=$(ondod keys add minter --home .ondo --keyring-backend test --output json | jq .address)
  ondod genesis add-genesis-account minter 10000000uusdc --home .ondo --keyring-backend test
  PAUSER=$(ondod keys add pauser --home .ondo --keyring-backend test --output json | jq .address)
  ondod genesis add-genesis-account pauser 10000000uusdc --home .ondo --keyring-backend test
  BLOCKLIST_OWNER=$(ondod keys add blocklist-owner --home .ondo --keyring-backend test --output json | jq .address)
  ondod genesis add-genesis-account blocklist-owner 10000000uusdc --home .ondo --keyring-backend test
  BLOCKLIST_PENDING_OWNER=$(ondod keys add blocklist-pending-owner --home .ondo --keyring-backend test --output json | jq .address)
  ondod genesis add-genesis-account blocklist-pending-owner 10000000uusdc --home .ondo --keyring-backend test
  ondod keys add user --home .ondo --keyring-backend test &> /dev/null
  ondod genesis add-genesis-account user 10000000uusdc --home .ondo --keyring-backend test

  TEMP=.ondo/genesis.json
  touch $TEMP && jq '.app_state.staking.params.bond_denom = "ustake"' .ondo/config/genesis.json > $TEMP && mv $TEMP .ondo/config/genesis.json
  touch $TEMP && jq '.app_state.usdy.blocklist_state.owner = '$BLOCKLIST_OWNER'' .ondo/config/genesis.json > $TEMP && mv $TEMP .ondo/config/genesis.json
  touch $TEMP && jq '.app_state.usdy.burner = '$BURNER'' .ondo/config/genesis.json > $TEMP && mv $TEMP .ondo/config/genesis.json
  touch $TEMP && jq '.app_state.usdy.minter = '$MINTER'' .ondo/config/genesis.json > $TEMP && mv $TEMP .ondo/config/genesis.json
  touch $TEMP && jq '.app_state.usdy.pauser = '$PAUSER'' .ondo/config/genesis.json > $TEMP && mv $TEMP .ondo/config/genesis.json

  ondod genesis gentx validator 1000000ustake --chain-id "ondo-1" --home .ondo --keyring-backend test &> /dev/null
  ondod genesis collect-gentxs --home .ondo &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .ondo/config/config.toml
fi

ondod start --home .ondo
