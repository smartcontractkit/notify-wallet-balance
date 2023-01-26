package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/go-resty/resty/v2"
)

// notifyStart sends a slack message to notify when the program starts
func notifyStart() error {
	restClient := resty.New()
	netsArr := []string{}
	for _, netConf := range GlobalConfig.NetworkConfigs {
		addrs := []string{fmt.Sprintf("*%s*\n", netConf.Name)}
		for _, addr := range netConf.Addresses {
			addrs = append(addrs, fmt.Sprintf("• %s", addr))
		}
		addrsStr := fmt.Sprintf(notifyStartNetwork, strings.Join(addrs, "\n"))
		netsArr = append(netsArr, addrsStr)
	}
	payload := fmt.Sprintf(notifyStartPayload, GlobalConfig.SlackChannel, GlobalConfig.SlackUser, strings.Join(netsArr, ",\n"))

	_, err := restClient.R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(payload).
		SetAuthToken(GlobalConfig.SlackAPIKey).
		Post(slackAPIURL)

	return err
}

// notifyStop sends a slack message notifying when the program has stopped running
func notifyStop() error {
	restClient := resty.New()

	payload := fmt.Sprintf(notifyStopPayload, GlobalConfig.SlackChannel, GlobalConfig.SlackUser)
	_, err := restClient.R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(payload).
		SetAuthToken(GlobalConfig.SlackAPIKey).
		Post(slackAPIURL)
	return err
}

// notifyAddress sends a slack message notifying when an address is low
func notifyAddress(network, address string, balance, limit *big.Float) error {
	restClient := resty.New()

	payload := fmt.Sprintf(
		notifyAddressPayload,
		GlobalConfig.SlackChannel,
		GlobalConfig.SlackUser,
		network,
		address,
		balance.String(),
		limit.String(),
	)
	_, err := restClient.R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(payload).
		SetAuthToken(GlobalConfig.SlackAPIKey).
		Post(slackAPIURL)
	return err
}

const (
	slackAPIURL        = "https://slack.com/api/chat.postMessage"
	notifyStartPayload = `{
	"channel": "%s",
	"blocks": [
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": ":white_check_mark: Started Monitoring Addresses :white_check_mark:",
				"emoji": true
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Notifying <@%s>"
			}
		},
		{
			"type": "divider"
		},
		%s
	]
}`
	notifyStartNetwork = `{
	"type": "section",
	"text": {
		"type": "mrkdwn",
		"text": "%s"
	}
},
{
	"type": "divider"
}`

	notifyStopPayload = `{
	"channel": "%s",
	"blocks": [
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": ":x: Monitoring Stopped :x:",
				"emoji": true
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Notifying <@%s>"
			}
		}
	]
}`
	notifyAddressPayload = `{
		"channel": "%s",
		"blocks": [
			{
				"type": "header",
				"text": {
					"type": "plain_text",
					"text": ":warning: Found Under-Funded Address :warning:",
					"emoji": true
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "Notifying <@%s>"
				}
			},
			{
				"type": "section",
				"text": {
					"type": "mrkdwn",
					"text": "Network: %s\nAddress: %s\nBalance: %s ETH\nLimit: %s"
				}
			}
		]
	}`
)
