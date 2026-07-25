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
	"strconv"

	listpkg "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/selectionv2"
)

func Select(projectName, policyName string) (int64, error) {
	model := selectionv2.NewModel(
		"Preheat Execution",
		fmt.Sprintf("Loading preheat executions for %s/%s...", projectName, policyName),
		loadItems(projectName, policyName),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return 0, fmt.Errorf("error during execution selection: %w", err)
	}

	selectionModel, ok := finalModel.(selectionv2.Model)
	if !ok {
		return 0, errors.New("unexpected execution selection result")
	}
	if selectionModel.Err != nil {
		return 0, selectionModel.Err
	}
	if selectionModel.Aborted {
		return 0, errors.New("user aborted execution selection")
	}
	if selectionModel.Choice == "" {
		return 0, errors.New("no execution selected")
	}

	id, err := strconv.ParseInt(selectionModel.Choice, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse execution ID %q: %v", selectionModel.Choice, err)
	}

	return id, nil
}

func loadItems(projectName, policyName string) selectionv2.Loader {
	return func() ([]listpkg.Item, error) {
		response, err := api.ListPreheatExecutions(projectName, policyName)
		if err != nil {
			return nil, err
		}
		if len(response.Payload) == 0 {
			return nil, errors.New("no preheat executions found")
		}

		items := make([]listpkg.Item, len(response.Payload))
		for i, execution := range response.Payload {
			items[i] = selectionv2.Item(strconv.FormatInt(execution.ID, 10))
		}

		return items, nil
	}
}
