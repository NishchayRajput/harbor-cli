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

	listpkg "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/selectionv2"
)

func Select(projectName string) (string, error) {
	model := selectionv2.NewModel(
		"Preheat Policy",
		fmt.Sprintf("Loading preheat policies for %s...", projectName),
		loadItems(projectName),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return "", fmt.Errorf("error during policy selection: %w", err)
	}

	selectionModel, ok := finalModel.(selectionv2.Model)
	if !ok {
		return "", errors.New("unexpected policy selection result")
	}
	if selectionModel.Err != nil {
		return "", selectionModel.Err
	}
	if selectionModel.Aborted {
		return "", errors.New("user aborted policy selection")
	}
	if selectionModel.Choice == "" {
		return "", errors.New("no policy selected")
	}

	return selectionModel.Choice, nil
}

func loadItems(projectName string) selectionv2.Loader {
	return func() ([]listpkg.Item, error) {
		response, err := api.ListPreheatPolicies(projectName, false)
		if err != nil {
			return nil, err
		}
		if len(response.Payload) == 0 {
			return nil, errors.New("no preheat policies found")
		}

		items := make([]listpkg.Item, len(response.Payload))
		for i, policy := range response.Payload {
			items[i] = selectionv2.Item(policy.Name)
		}

		return items, nil
	}
}
