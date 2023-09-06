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

// Initialize logging
func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// Starts monitoring and waits for issues
func main() {
	log.Info().Msg("Starting")
	err := loadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}

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
	if err = notifyStart(); err != nil {
		mainErrChan <- err
	}

	signal.Notify(terminationChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case mainErr := <-mainErrChan:
			if err = notifyStop(); err != nil {
				log.Error().Err(err).Msg("Error while trying to notify of stopped runner")
			}
			log.Fatal().Err(mainErr).Msg("Unrecoverable Error Monitoring Chain")
		case <-terminationChan:
			if err = notifyStop(); err != nil {
				log.Error().Err(err).Msg("Error while trying to notify of stopped runner")
			}
			log.Fatal().Msg("Monitoring Killed!")
		}
	}
}

// monitorNetwork polls addresses based on the network's poll interval
func monitorNetwork(netConf *NetworkConfig, mainErrChan chan error) {
	log.Info().
		Interface("Addresses", netConf.Addresses).
		Str("URL", netConf.URL).
		Str("Network", netConf.Name).
		Msg("Monitoring Network")
	pollInterval := time.NewTicker(GlobalConfig.NotificationInterval)

	for range pollInterval.C {
		log.Info().Str("Network", netConf.Name).Msg("Checking Addresses")
		client, err := ethclient.Dial(netConf.URL)
		if err != nil {
			mainErrChan <- err
		}
		for _, address := range netConf.Addresses {
			err = checkAddress(client, netConf, address)
			if err != nil {
				mainErrChan <- err
			}
		}
	}
}

// checks a provided address with a provided client once, notifying if the balance is too low
func checkAddress(client *ethclient.Client, netConf *NetworkConfig, addressString string) error {
	address := common.HexToAddress(addressString)
	bigBal, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return err
	}
	balance := weiToEther(bigBal)
	bigLowerLimit := big.NewFloat(netConf.LowerLimit * 1.0)
	if balance.Cmp(bigLowerLimit) <= 0 {
		log.Warn().
			Str("Lower Limit", bigLowerLimit.String()).
			Str("Balance", balance.String()).
			Str("Network", netConf.Name).
			Msg("Address Below Limit!")

		if err = notifyAddress(netConf, addressString, balance, bigLowerLimit); err != nil {
			log.Error().Err(err).Msg("Error trying to notify of under-funded address")
		}

	} else {
		log.Debug().
			Str("Lower Limit", bigLowerLimit.String()).
			Str("Balance", balance.String()).
			Msg("Address Balance Fine")
	}

	return nil
}

// weiToEther gives a quick conversion of wei units to ETH for better readability
func weiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}
