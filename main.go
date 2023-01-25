package main

import (
	"context"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	EthMult uint64 = 1e18
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	log.Info().Msg("Starting")
	err := loadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}
	log.Info().Msg("Loaded Config")

	terminationChan := make(chan os.Signal, 1)
	mainErrChan := make(chan error, 1)
	for _, networkConf := range GlobalConfig.NetworkConfigs {
		c, err := ethclient.Dial(networkConf.URL)
		if err != nil { // Validate all network URLs on initialization
			log.Fatal().
				Err(err).
				Str("URL", networkConf.URL).
				Str("Network", networkConf.Name).
				Msg("Error on initially connecting to network")
		}
		c.Close()
		go monitorNetwork(networkConf, mainErrChan)
	}

	signal.Notify(terminationChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case err = <-mainErrChan:
			log.Fatal().Err(err).Msg("Unrecoverable Error Monitoring Chain")
		case <-terminationChan:
			log.Fatal().Msg("Monitoring Killed!")
		}
	}
}

// monitorNetwork polls addresses based on the network's poll interval
func monitorNetwork(netConf *NetworkConfig, mainErrChan chan error) {
	log.Info().
		Str("Poll Interval", netConf.PollInterval.String()).
		Interface("Addresses", netConf.Addresses).
		Str("URL", netConf.URL).
		Str("Network", netConf.Name).
		Msg("Monitoring Network")
	pollInterval := time.NewTicker(netConf.PollInterval)

	for {
		select {
		case <-pollInterval.C:
			log.Info().Str("Network", netConf.Name).Msg("Checking Addresses")
			client, err := ethclient.Dial(netConf.URL)
			if err != nil {
				mainErrChan <- err
			}
			for _, address := range netConf.Addresses {
				err = checkAddress(client, address, netConf.LowerLimit)
				if err != nil {
					mainErrChan <- err
				}
			}
		}
	}
}

// checks a provided address with a provided client once, notifying if the balance is too low
func checkAddress(client *ethclient.Client, addressString string, lowerLimit float64) error {
	address := common.HexToAddress(addressString)
	bigBal, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return err
	}
	balance := weiToEther(bigBal)
	bigLowerLimit := big.NewFloat(lowerLimit * 1.0)
	if balance.Cmp(bigLowerLimit) <= 0 {
		log.Warn().
			Str("Lower Limit", bigLowerLimit.String()).
			Str("Balance", balance.String()).
			Msg("Address Below Limit!")
	} else {
		log.Debug().
			Str("Lower Limit", bigLowerLimit.String()).
			Str("Balance", balance.String()).
			Msg("Address Balance Fine")
	}

	return nil
}

func weiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}
