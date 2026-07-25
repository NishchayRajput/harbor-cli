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

package execution

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

var viewColumns = []table.Column{
	{Title: "ID", Width: tablelist.WidthS},
	{Title: "Status", Width: tablelist.WidthM},
	{Title: "Trigger", Width: tablelist.WidthM},
	{Title: "Success Rate", Width: tablelist.WidthM},
	{Title: "Start Time", Width: tablelist.WidthL},
	{Title: "End Time", Width: tablelist.WidthL},
	{Title: "Vendor", Width: tablelist.WidthM},
}

func View(projectName, policyName string, executionID int64) error {
	model := tablelistv2.NewModel(
		viewColumns,
		fmt.Sprintf("Loading preheat execution %d for %s/%s...", executionID, projectName, policyName),
		loadViewRows(projectName, policyName, executionID),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running preheat execution view: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected preheat execution view model result")
	}

	return loadedModel.Err
}

func loadViewRows(projectName, policyName string, executionID int64) tablelistv2.Loader {
	return func() ([]table.Row, error) {
		response, err := api.GetPreheatExecution(projectName, policyName, executionID)
		if err != nil {
			return nil, fmt.Errorf("failed to view preheat execution: %v", utils.ParseHarborErrorMsg(err))
		}

		return []table.Row{buildRow(response.Payload)}, nil
	}
}
