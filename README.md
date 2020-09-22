# Fork
Fork of [tendermint-validator](https://gitlab.com/polychainlabs/tendermint-validator) and ongoing work on tendermint threshold validator using [threshold-ed25519](https://gitlab.com/polychainlabs/threshold-ed25519)

Read article [Threshold Validator for Tendermint
](https://blog.polychainlabs.com/tendermint/2020/03/26/threshold-validator-for-tendermint.html) by Polychain Labs

With respect to Polychain Labs developers ([Roman Shtylman](https://github.com/defunctzombie)) and [Cybernetic Destiny](https://www.mintscan.io/validators/cosmosvaloper1d7ufwp2rgfj7s7pfw2q7vm2lc9txmr8vh77ztr) validator

Updated to work with tendermint v0.33

# Tendermint Validator

A lightweight single key tendermint validator for sentry nodes.

## Design

A lightweight alternative to using a full node instance for validating blocks. The validator is able to connect to any number of sentry nodes and will sign blocks provided by the nodes. The validator maintains a watermark file to protect against double signing.

## Pre-requisites

Before starting, please make sure to fully understand node and validator requirements and operation for your particular network and chain.

## Setup

_The security of any key material is outside the scope of this guide. At a minimum we recommend performing key material steps on airgapped computers and using your audited security procedures._

### Setup Validator Instance

Configure the instance with a [toml](https://github.com/toml-lang/toml) file. Below is a sample configuration.

```toml
# Path to priv validator key json file
key_file = "/path/to/priv_validator_key.json"

# The state directory stores watermarks for double signing protection.
# Each validator instance maintains a watermark.
state_dir = "/path/to/state/dir"

# The network chain id for your p2p nodes
chain_id = "chain-id-here"

# Configure any number of p2p network nodes.
# We recommend at least 2 nodes for redundancy.
[[node]]
address = "tcp://<node-a ip>:1234"

[[node]]
address = "tcp://<node-b ip>:1234"
```

## Configure p2p network nodes

Validators are not directly connected to the p2p network nor do they store chain and application state. They rely on nodes to receive blocks from the p2p network, make signing requests, and relay the signed blocks back to the p2p network.

To make a node available as a relay for a validator, find the `priv_validator_laddr` (or equivalent) configuration item in your node's configuration file. Change this value, to accept connections on an IP address and port of your choosing.

```diff
 # TCP or UNIX socket address for Tendermint to listen on for
 # connections from an external PrivValidator process
-priv_validator_laddr = ""
+priv_validator_laddr = "tcp://0.0.0.0:1234"
```

_Full configuration and operation of your tendermint node is outside the scope of this guide. You should consult your network's documentation for node configuration._

_We recommend hosting nodes on separate and isolated infrastructure from your validator instances._

## Launch validator

Once your validator instance and node is configured, you can launch the signer.

```bash
signer --config /path/to/config.toml
```

_We recommend using systemd or similar service management program as appropriate for your runtime platform._

## Security

Security and management of any key material is outside the scope of this service. Always consider your own security and risk profile when dealing with sensitive keys, services, or infrastructure.

## No Liability

As far as the law allows, this software comes as is,
without any warranty or condition, and no contributor
will be liable to anyone for any damages related to this
software or this license, under any kind of legal claim.

## References

- https://docs.tendermint.com/master/tendermint-core/validators.html
- https://hub.cosmos.network/master/validators/overview.html
