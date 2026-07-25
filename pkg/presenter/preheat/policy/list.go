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

package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelistv2"
)

var listColumns = []table.Column{
	{Title: "Name", Width: tablelist.WidthM},
	{Title: "Enabled", Width: tablelist.WidthS},
	{Title: "Provider", Width: tablelist.WidthXL},
	{Title: "Filters", Width: tablelist.Width3XL},
	{Title: "Trigger", Width: tablelist.WidthM},
	{Title: "Creation Time", Width: tablelist.WidthL},
	{Title: "Description", Width: tablelist.WidthM},
}

func List(projectName string, opts api.ListFlags) error {
	model := tablelistv2.NewModel(
		listColumns,
		fmt.Sprintf("Loading preheat policies for %s...", projectName),
		loadRows(projectName, opts),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running preheat policy list: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected preheat policy list model result")
	}

	return loadedModel.Err
}

func loadRows(projectName string, opts api.ListFlags) tablelistv2.Loader {
	return func() ([]table.Row, error) {
		response, err := api.ListPreheatPolicies(projectName, false, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list preheat policies: %v", utils.ParseHarborErrorMsg(err))
		}

		return buildRows(response.Payload), nil
	}
}

func buildRows(policies []*models.PreheatPolicy) []table.Row {
	rows := make([]table.Row, 0, len(policies))

	for _, policy := range policies {
		enabled := "No"
		if policy.Enabled {
			enabled = "Yes"
		}

		createdTime, _ := utils.FormatCreatedTime(policy.CreationTime.String())
		rows = append(rows, table.Row{
			policy.Name,
			enabled,
			policy.ProviderName,
			formatFilters(policy.Filters),
			formatTrigger(policy.Trigger),
			createdTime,
			policy.Description,
		})
	}

	return rows
}

func formatFilters(raw string) string {
	var filters []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
	var filterParts []string

	if err := json.Unmarshal([]byte(raw), &filters); err == nil {
		for _, filter := range filters {
			filterParts = append(filterParts, fmt.Sprintf("%s: %s", filter.Type, filter.Value))
		}
	}

	return strings.Join(filterParts, "  ")
}

func formatTrigger(raw string) string {
	var trigger struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal([]byte(raw), &trigger); err != nil {
		return raw
	}

	return trigger.Type
}
