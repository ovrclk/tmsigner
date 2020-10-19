module github.com/ovrclk/tmsigner

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/cosmos/cosmos-sdk v0.40.0-rc0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	github.com/tendermint/tendermint v0.34.0-rc4.0.20201005135527-d7d0ffea13c6
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
