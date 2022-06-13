module github.com/Gravity-Bridge/Gravity-Bridge/module

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.45.5-0.20220523154235-2921a1c3c918
	github.com/cosmos/ibc-go/v2 v2.1.0
	github.com/ethereum/go-ethereum v1.10.16
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/osmosis-labs/bech32-ibc v0.2.0-rc2
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/rs/zerolog v1.25.0 // indirect
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.12.0
	github.com/stretchr/testify v1.7.2
	github.com/tendermint/tendermint v0.34.20-0.20220517115723-e6f071164839
	github.com/tendermint/tm-db v0.6.7
	github.com/tharsis/ethermint v0.16.1
	google.golang.org/genproto v0.0.0-20220602131408-e326c6e8e9c8
	google.golang.org/grpc v1.46.2
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/gogo/grpc => google.golang.org/grpc v1.33.2

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2
