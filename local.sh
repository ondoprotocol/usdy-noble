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
  ondod keys add user --home .ondo --keyring-backend test &> /dev/null
  ondod genesis add-genesis-account user 10000000uusdc --home .ondo --keyring-backend test

  TEMP=.ondo/genesis.json
  touch $TEMP && jq '.app_state.staking.params.bond_denom = "ustake"' .ondo/config/genesis.json > $TEMP && mv $TEMP .ondo/config/genesis.json

  ondod genesis gentx validator 1000000ustake --chain-id "ondo-1" --home .ondo --keyring-backend test &> /dev/null
  ondod genesis collect-gentxs --home .ondo &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .ondo/config/config.toml
fi

ondod start --home .ondo
