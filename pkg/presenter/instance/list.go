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

package instance

import (
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelistv2"
)

var listColumns = []table.Column{
	{Title: "ID", Width: tablelist.WidthXS},
	{Title: "Name", Width: tablelist.WidthM},
	{Title: "Provider", Width: tablelist.WidthM},
	{Title: "Endpoint", Width: tablelist.WidthXXL},
	{Title: "Status", Width: tablelist.WidthM},
	{Title: "Auth Mode", Width: tablelist.WidthM},
	{Title: "Description", Width: tablelist.WidthM},
	{Title: "Default", Width: tablelist.WidthS},
	{Title: "Insecure", Width: tablelist.WidthS},
	{Title: "Enabled", Width: tablelist.WidthS},
	{Title: "Setup Timestamp", Width: tablelist.WidthXL},
}

func List(opts api.ListFlags) error {
	model := tablelistv2.NewModel(
		listColumns,
		"Loading instances...",
		loadRows(opts),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running instance list: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected instance list model result")
	}

	return loadedModel.Err
}

func loadRows(opts api.ListFlags) tablelistv2.Loader {
	return func() ([]table.Row, error) {
		response, err := api.ListAllInstance(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get instance list: %v", err)
		}
		if len(response.Payload) == 0 {
			return nil, errors.New("no instances found")
		}

		return buildRows(response.Payload), nil
	}
}

func buildRows(instances []*models.Instance) []table.Row {
	rows := make([]table.Row, 0, len(instances))
	for _, instance := range instances {
		rows = append(rows, buildRow(instance))
	}
	return rows
}

func buildRow(instance *models.Instance) table.Row {
	return table.Row{
		fmt.Sprintf("%d", instance.ID),
		instance.Name,
		instance.Vendor,
		instance.Endpoint,
		instance.Status,
		instance.AuthMode,
		instance.Description,
		fmt.Sprintf("%t", instance.Default),
		fmt.Sprintf("%t", instance.Insecure),
		fmt.Sprintf("%t", instance.Enabled),
		time.Unix(instance.SetupTimestamp, 0).Format("2006-01-02 15:04:05"),
	}
}
