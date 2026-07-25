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
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelistv2"
)

var viewColumns = []table.Column{
	{Title: "Name", Width: tablelist.WidthM},
	{Title: "Enabled", Width: tablelist.WidthS},
	{Title: "Provider", Width: tablelist.WidthXL},
	{Title: "Filters", Width: tablelist.Width3XL},
	{Title: "Trigger", Width: tablelist.WidthM},
	{Title: "Creation Time", Width: tablelist.WidthL},
	{Title: "Updated", Width: tablelist.WidthL},
	{Title: "Description", Width: tablelist.WidthM},
}

func View(projectName, policyName string) error {
	model := tablelistv2.NewModel(
		viewColumns,
		fmt.Sprintf("Loading preheat policy %s for %s...", policyName, projectName),
		loadViewRows(projectName, policyName),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running preheat policy view: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected preheat policy view model result")
	}

	return loadedModel.Err
}

func loadViewRows(projectName, policyName string) tablelistv2.Loader {
	return func() ([]table.Row, error) {
		response, err := api.GetPreheatPolicy(projectName, policyName)
		if err != nil {
			return nil, fmt.Errorf("failed to get preheat policy: %v", utils.ParseHarborErrorMsg(err))
		}

		return buildViewRows(response.Payload), nil
	}
}

func buildViewRows(policy *models.PreheatPolicy) []table.Row {
	enabled := "No"
	if policy.Enabled {
		enabled = "Yes"
	}

	createdTime, _ := utils.FormatCreatedTime(policy.CreationTime.String())
	updatedTime, _ := utils.FormatCreatedTime(policy.UpdateTime.String())

	return []table.Row{{
		policy.Name,
		enabled,
		policy.ProviderName,
		formatFilters(policy.Filters),
		formatTrigger(policy.Trigger),
		createdTime,
		updatedTime,
		policy.Description,
	}}
}
