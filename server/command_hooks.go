package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	cells_sdk "github.com/pydio/cells-sdk-go"
	"github.com/pydio/cells-sdk-go/client/tree_service"
	"github.com/pydio/cells-sdk-go/example/cmd"
	"github.com/pydio/cells-sdk-go/models"
)

const (
	commandTriggerList = "list"
)

func (p *Plugin) registerCommands() error {
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          commandTriggerList,
		AutoComplete:     true,
		AutoCompleteDesc: "Lists files on Cells",
		AutoCompleteHint: "(something)",
	}); err != nil {
		return err
	}
	return nil
}

func (p *Plugin) ExecuteCommand(ctx *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")
	switch trigger {
	case commandTriggerList:
		return p.executeListCommand(args), nil
	default:
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknown command: " + args.Command),
		}, nil
	}
}

func (p *Plugin) executeListCommand(args *model.CommandArgs) *model.CommandResponse {
	config := &cells_sdk.SdkConfig{
		Url:        "https://cells-test.your-files-your-rules.eu",
		ClientKey:  "cells-front",
		User:       "",
		Password:   "",
		SkipVerify: false,
	}

	ctx, cli, err := cmd.GetApiClient(config)
	if err != nil {
		errorMessage := "failed to get client"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         errorMessage,
		}
	}

	var target string
	if strings.Fields(args.Command)[1] == "" {
		errorMessage := "no target provided"
		p.API.LogError(errorMessage)
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         errorMessage,
		}
	} else {
		target = strings.Fields(args.Command)[1]
	}
	params := &tree_service.BulkStatNodesParams{
		Context: ctx,
		Body:    &models.RestGetBulkMetaRequest{NodePaths: []string{target}},
	}

	resp, err := cli.TreeService.BulkStatNodes(params)
	if err != nil {
		errorMessage := "failed to get bulk"
		p.API.LogError(errorMessage, "err", err.Error())
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         errorMessage,
		}
	}
	var aa []*model.SlackAttachment
	for _, node := range resp.Payload.Nodes {
		aa = append(aa, &model.SlackAttachment{
			Title: node.MetaStore["name"],
			Text:  node.Path,
		})
	}

	post := &model.Post{
		ChannelId: args.ChannelId,
		Message:   fmt.Sprintf("found nodes %+v", resp.Payload.Nodes),
		Props: model.StringInterface{
			"attachments": aa,
		},
	}

	_ = p.API.SendEphemeralPost(args.UserId, post)

	return &model.CommandResponse{}
}
