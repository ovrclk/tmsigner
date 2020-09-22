#!/bin/bash

SIGNER_DATA=$HOME/.tmsigner
CHAINID=$2

# Ensure user understands what will be deleted
if [[ -d $SIGNER_DATA ]] && [[ ! "$1" == "skip" ]]; then
  read -p "$0 will delete \$HOME/.tmsigner folder. Do you wish to continue? (y/n): " -n 1 -r
  echo 
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      exit 1
  fi
fi

if [ -z "$2" ]; then
  echo "Need to input chain-id..."
  exit 1
fi

rm -rf $SIGNER_DATA &> /dev/null
mkdir -p $SIGNER_DATA/data
touch $SIGNER_DATA/config.toml

# Path to priv validator key json file
echo "key_file = \"$SIGNER_DATA/priv_validator_key.json\"" >> $SIGNER_DATA/config.toml
echo "state_dir = \"$SIGNER_DATA/data\"" >> $SIGNER_DATA/config.toml
echo "chain_id = \"$CHAINID\"" >> $SIGNER_DATA/config.toml
echo "[[ node ]]" >> $SIGNER_DATA/config.toml
echo "address = \"tcp://localhost:1234\"" >> $SIGNER_DATA/config.toml
