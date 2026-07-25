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

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelistv2"
)

var viewColumns = listColumns

func View(instanceName string, isID bool) error {
	model := tablelistv2.NewModel(
		viewColumns,
		fmt.Sprintf("Loading instance %s...", instanceName),
		loadViewRows(instanceName, isID),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running instance view: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected instance view model result")
	}

	return loadedModel.Err
}

func loadViewRows(instanceName string, isID bool) tablelistv2.Loader {
	return func() ([]table.Row, error) {
		response, err := api.GetInstance(instanceName, isID)
		if err != nil {
			return nil, fmt.Errorf("failed to get instance: %v", utils.ParseHarborErrorMsg(err))
		}
		if response == nil || response.Payload == nil {
			return nil, errors.New("failed to get instance: empty response")
		}

		return []table.Row{buildRow(response.Payload)}, nil
	}
}
