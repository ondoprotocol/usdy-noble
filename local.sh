alias aurad=./simapp/build/simd

for arg in "$@"
do
    case $arg in
        -r|--reset)
        rm -rf .aura
        shift
        ;;
    esac
done

if ! [ -f .aura/data/priv_validator_state.json ]; then
  aurad init validator --chain-id "aura-1" --home .aura &> /dev/null

  aurad keys add validator --home .aura --keyring-backend test &> /dev/null
  aurad add-genesis-account validator 1000000ustake --home .aura --keyring-backend test
  OWNER=$(aurad keys add owner --home .aura --keyring-backend test --output json | jq .address)
  aurad add-genesis-account owner 10000000uusdc --home .aura --keyring-backend test
  PENDING_OWNER=$(aurad keys add pending-owner --home .aura --keyring-backend test --output json | jq .address)
  aurad add-genesis-account pending-owner 10000000uusdc --home .aura --keyring-backend test
  BURNER=$(aurad keys add burner --home .aura --keyring-backend test --output json | jq .address)
  aurad add-genesis-account burner 10000000uusdc --home .aura --keyring-backend test
  MINTER=$(aurad keys add minter --home .aura --keyring-backend test --output json | jq .address)
  aurad add-genesis-account minter 10000000uusdc --home .aura --keyring-backend test
  PAUSER=$(aurad keys add pauser --home .aura --keyring-backend test --output json | jq .address)
  aurad add-genesis-account pauser 10000000uusdc --home .aura --keyring-backend test
  BLOCKLIST_OWNER=$(aurad keys add blocklist-owner --home .aura --keyring-backend test --output json | jq .address)
  aurad add-genesis-account blocklist-owner 10000000uusdc --home .aura --keyring-backend test
  BLOCKLIST_PENDING_OWNER=$(aurad keys add blocklist-pending-owner --home .aura --keyring-backend test --output json | jq .address)
  aurad add-genesis-account blocklist-pending-owner 10000000uusdc --home .aura --keyring-backend test
  aurad keys add alice --home .aura --keyring-backend test &> /dev/null
  aurad add-genesis-account alice 10000000uusdc --home .aura --keyring-backend test
  aurad keys add bob --home .aura --keyring-backend test &> /dev/null
  aurad add-genesis-account bob 10000000uusdc --home .aura --keyring-backend test
  aurad keys add charlie --home .aura --keyring-backend test &> /dev/null
  aurad add-genesis-account charlie 10000000uusdc --home .aura --keyring-backend test

  TEMP=.aura/genesis.json
  touch $TEMP && jq '.app_state.staking.params.bond_denom = "ustake"' .aura/config/genesis.json > $TEMP && mv $TEMP .aura/config/genesis.json
  touch $TEMP && jq '.app_state.aura.blocklist_state.owner = '$BLOCKLIST_OWNER'' .aura/config/genesis.json > $TEMP && mv $TEMP .aura/config/genesis.json
  touch $TEMP && jq '.app_state.aura.owner = '$OWNER'' .aura/config/genesis.json > $TEMP && mv $TEMP .aura/config/genesis.json
  touch $TEMP && jq '.app_state.aura.burners = [{"address": '$BURNER', "allowance": "0"}]' .aura/config/genesis.json > $TEMP && mv $TEMP .aura/config/genesis.json
  touch $TEMP && jq '.app_state.aura.minters = [{"address": '$MINTER', "allowance": "100000000000000000000"}]' .aura/config/genesis.json > $TEMP && mv $TEMP .aura/config/genesis.json
  touch $TEMP && jq '.app_state.aura.pausers = ['$PAUSER']' .aura/config/genesis.json > $TEMP && mv $TEMP .aura/config/genesis.json

  aurad gentx validator 1000000ustake --chain-id "aura-1" --home .aura --keyring-backend test &> /dev/null
  aurad collect-gentxs --home .aura &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .aura/config/config.toml
fi

aurad start --home .aura
