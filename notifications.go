package main

import (
	"fmt"
	"math/big"

	"github.com/slack-go/slack"
)

// notifyStart sends a slack message to notify when the program starts
func notifyStart() error {
	notificationBlocks := []slack.Block{}
	notificationBlocks = append(notificationBlocks, slack.NewTextBlockObject("mrkdwn", "# Started Monitoring Addresses", true, false))
	if len(GlobalConfig.SlackUsers) > 0 {
		userNoti := "Notifying"
		for _, user := range GlobalConfig.SlackUsers {
			userNoti = fmt.Sprintf("%s <@%s>", userNoti, user)
		}
		notificationBlocks = append(notificationBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", userNoti, false, true), nil, nil))
	}

	notificationBlocks = append(notificationBlocks, slack.NewDividerBlock())

	for _, netConf := range GlobalConfig.NetworkConfigs {
		netText := fmt.Sprintf("### %s\nLower Limit: %f", netConf.Name, netConf.LowerLimit)
		for _, addr := range netConf.Addresses {
			netText = fmt.Sprintf("%s\n* %s\n", netText, addr)
		}
		notificationBlocks = append(notificationBlocks, slack.NewTextBlockObject("mrkdwn", netText, true, false))
	}

	slackClient := slack.New(GlobalConfig.SlackAPIKey)
	msgOptionBlocks := []slack.MsgOption{slack.MsgOptionBlocks(notificationBlocks...), slack.MsgOptionAsUser(true)}
	_, _, err := slackClient.PostMessage(GlobalConfig.SlackChannel, msgOptionBlocks...)
	return err
}

// notifyStop sends a slack message notifying when the program has stopped running
func notifyStop() error {
	notificationBlocks := []slack.Block{}
	notificationBlocks = append(notificationBlocks, slack.NewTextBlockObject("mrkdwn", "# :x: Stopped Monitoring Addresses! :x:", true, false))
	if len(GlobalConfig.SlackUsers) > 0 {
		userNoti := "Notifying"
		for _, user := range GlobalConfig.SlackUsers {
			userNoti = fmt.Sprintf("%s <@%s>", userNoti, user)
		}
		notificationBlocks = append(notificationBlocks, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", userNoti, false, true), nil, nil))
	}

	slackClient := slack.New(GlobalConfig.SlackAPIKey)
	msgOptionBlocks := []slack.MsgOption{slack.MsgOptionBlocks(notificationBlocks...), slack.MsgOptionAsUser(true)}
	_, _, err := slackClient.PostMessage(GlobalConfig.SlackChannel, msgOptionBlocks...)
	return err
}

// notifyAddress sends a slack message notifying when an address is low
func notifyAddress(address string, balance, limit *big.Float) error {
	notificationBlocks := []slack.Block{}
	notificationBlocks = append(notificationBlocks,
		slack.NewTextBlockObject("mrkdwn",
			fmt.Sprintf("# :warning: Address Under-Funded :warning:\n\nAddress: %s | Balance: %s | Lower Limit: %s",
				address, balance.String(), limit.String(),
			),
			true, false,
		),
	)

	slackClient := slack.New(GlobalConfig.SlackAPIKey)
	msgOptionBlocks := []slack.MsgOption{slack.MsgOptionBlocks(notificationBlocks...), slack.MsgOptionAsUser(true)}
	_, _, err := slackClient.PostMessage(GlobalConfig.SlackChannel, msgOptionBlocks...)
	return err
}
