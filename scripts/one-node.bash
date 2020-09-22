#!/bin/sh
# USAGE: ./one-chain test-chain-id ./data 26657 26656

CHAINID=$1
CHAINDIR=$2

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input directory to create files in..."
  exit 1
fi

echo "Creating akashd instance with home=$CHAINDIR chain-id=$CHAINID..."
# Build genesis file incl account for passed address
coins="100000000000stake,100000000000samoleans"
akashd --home $CHAINDIR/$CHAINID --chain-id $CHAINID init $CHAINID
akashctl --home $CHAINDIR/$CHAINID keys add validator --keyring-backend="test"
akashd --home $CHAINDIR/$CHAINID add-genesis-account $(akashctl --home $CHAINDIR/$CHAINID keys --keyring-backend="test" show validator -a) $coins 
echo "ABOUT TO CREATE GENTX"
akashd --home $CHAINDIR/$CHAINID gentx --home-client $CHAINDIR/$CHAINID --name validator --keyring-backend="test"
akashd --home $CHAINDIR/$CHAINID collect-gentxs

# Set proper defaults and change ports
# TODO: sed for linux
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAINDIR/$CHAINID/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAINDIR/$CHAINID/config/config.toml
sed -i '' 's#priv_validator_laddr = ""#priv_validator_laddr = "tcp://0.0.0.0:1234"#g' $CHAINDIR/$CHAINID/config/config.toml

cp $CHAINDIR/$CHAINID/config/priv_validator_key.json $HOME/.tmsigner/priv_validator_key.json
cp $CHAINDIR/$CHAINID/data/priv_validator_state.json $HOME/.tmsigner/data/${CHAINID}_priv_validator_state.json

# Start the gaia
akashd --home $CHAINDIR/$CHAINID start --pruning=nothing > $CHAINDIR/$CHAINID.log 2>&1 &