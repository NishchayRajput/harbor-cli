// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package context

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelistv2"
)

var contextColumns = []table.Column{
	{Title: "Name", Width: tablelist.Width3XL},
	{Title: "Username", Width: tablelist.WidthL},
	{Title: "Server Address", Width: tablelist.WidthXXL},
}

func ListContexts() error {
	model := tablelistv2.NewModel(
		contextColumns,
		"Loading contexts...",
		loadContextRows(),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running context list: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected context list model result")
	}

	return loadedModel.Err
}

func loadContextRows() tablelistv2.Loader {
	return func() ([]table.Row, error) {
		config, err := utils.GetCurrentHarborConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get config: %w", err)
		}

		contexts := buildContextViews(config)
		if len(contexts) == 0 {
			return nil, errors.New("no contexts found")
		}

		return buildContextRows(contexts, config.CurrentCredentialName), nil
	}
}

func buildContextViews(config *utils.HarborConfig) []api.ContextListView {
	contexts := make([]api.ContextListView, 0, len(config.Credentials))

	for _, cred := range config.Credentials {
		contexts = append(contexts, api.ContextListView{
			Name:     cred.Name,
			Username: cred.Username,
			Server:   cred.ServerAddress,
		})
	}

	return contexts
}

func buildContextRows(contexts []api.ContextListView, currentCredential string) []table.Row {
	rows := make([]table.Row, 0, len(contexts))

	for _, ctx := range contexts {
		name := "  " + ctx.Name
		if ctx.Name == currentCredential {
			name = "* " + ctx.Name
		}

		rows = append(rows, table.Row{
			name,
			ctx.Username,
			ctx.Server,
		})
	}

	return rows
}
