package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
    "id": "com.mattermost.my-plugin",
    "name": "My Plugin",
    "description": "A plugin to enhance Mattermost.",
    "homepage_url": "https://github.com/mattermost/mattermost-plugin-starter-template",
    "support_url": "https://github.com/mattermost/mattermost-plugin-starter-template/issues",
    "release_notes_url": "https://github.com/mattermost/mattermost-plugin-starter-template/releases/tag/v0.1.0",
    "icon_path": "assets/starter-template-icon.svg",
    "version": "0.1.0",
    "min_server_version": "5.12.0",
    "server": {
        "executables": {
            "linux-amd64": "server/dist/plugin-linux-amd64",
            "darwin-amd64": "server/dist/plugin-darwin-amd64",
            "windows-amd64": "server/dist/plugin-windows-amd64.exe"
        }
    },
    "webapp": {
        "bundle_path": "webapp/dist/main.js"
    },
    "settings_schema": {
        "header": "",
        "footer": "",
        "settings": []
    }
}`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
