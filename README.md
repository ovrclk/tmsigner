# `tmsigner`

A lightweight single key tendermint validator signer for use with an array of sentry nodes.

### Fork information and code history

The `tmsigner` is based off of [litvintech](https://github.com/litvintech/tendermint-validator)'s fork of [polychainlabs/tendermint-validator](https://gitlab.com/polychainlabs/tendermint-validator). This fork reorganizes and adds much nicer config management to the code. This is the basis for future improvements to this code.

With respect to the work by:
- [Roman Shtylman](https://github.com/defunctzombie)
- [Cybernetic Destiny](https://www.mintscan.io/validators/cosmosvaloper1d7ufwp2rgfj7s7pfw2q7vm2lc9txmr8vh77ztr)
- [`litvintech`](https://github.com/litvintech)

### Design

A lightweight alternative to using a full node instance for validating blocks. The validator is able to connect to any number of sentry nodes and will sign blocks provided by the nodes. The validator maintains a watermark file to protect against double signing.

## Pre-requisites

Before starting, please make sure to fully understand node and validator requirements and operation for your particular network and chain.

## Setup

_The security of any key material is outside the scope of this guide. At a minimum we recommend performing key material steps on airgapped computers and using your audited security procedures._

### Setup Validator Instance

Configure the instance with a [toml](https://github.com/toml-lang/toml) file. Below is a sample configuration.

```toml
# The network chain id for your p2p nodes
chain_id = "chain-id-here"

# Configure any number of p2p network nodes.
# We recommend at least 2 nodes for redundancy.
[[node]]
address = "tcp://<node-a ip>:1234"

[[node]]
address = "tcp://<node-b ip>:1234"
```

You can generate this by running `tmsigner init {{chain_id}}`. `tmsigner` expects the private key file to exist in `$HOME/.tmsigner/priv_validator_key.json` and the state file to exist in `$HOME/.tmsigner/data/{{chain_id}}_priv_validator_state.json`. You can change the default home folder using the `--home` flag on all commands.

### Configure p2p network nodes

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
tmsigner start
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
- https://gitlab.com/polychainlabs/tendermint-validator
