package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"sync"
	"time"

	"github.com/ovrclk/tmsigner/signer"

	tmlog "github.com/tendermint/tendermint/libs/log"
	tos "github.com/tendermint/tendermint/libs/os"
	svc "github.com/tendermint/tendermint/libs/service"
	"github.com/tendermint/tendermint/privval"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	logger := tmlog.NewTMLogger(
		tmlog.NewSyncWriter(os.Stdout),
	).With("module", "validator")

	var configFile = flag.String("config", "", "path to configuration file")
	flag.Parse()

	if *configFile == "" {
		panic("--config flag is required")
	}

	config, err := signer.LoadConfigFromFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(
		"Tendermint Validator",
		"priv-key", config.PrivValKeyFile,
		"priv-state-dir", config.PrivValStateDir,
	)

	signer.InitSerialization()

	// services to stop on shutdown
	var services []svc.Service

	chainID := config.ChainID
	if chainID == "" {
		log.Fatal("chain_id option is required")
	}

	stateFile := path.Join(config.PrivValStateDir, fmt.Sprintf("%s_priv_validator_state.json", chainID))

	if !fileExists(stateFile) {
		log.Fatalf("State file missing: %s\n", stateFile)
	}

	val := privval.LoadFilePV(config.PrivValKeyFile, stateFile)
	pv := &signer.PvGuard{PrivValidator: val}

	for _, node := range config.Nodes {
		dialer := net.Dialer{Timeout: 30 * time.Second}
		signer := signer.NewNodeClient(node.Address, logger, config.ChainID, pv, dialer)

		err := signer.Start()
		if err != nil {
			panic(err)
		}

		services = append(services, signer)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	tos.TrapSignal(logger, func() {
		for _, service := range services {
			err := service.Stop()
			if err != nil {
				panic(err)
			}
		}
		wg.Done()
	})
	wg.Wait()
}
