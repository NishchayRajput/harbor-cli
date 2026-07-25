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

	listpkg "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/selectionv2"
)

func Select() (string, error) {
	model := selectionv2.NewModel(
		"Instance",
		"Loading instances...",
		loadItems(),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return "", fmt.Errorf("error during instance selection: %w", err)
	}

	selectionModel, ok := finalModel.(selectionv2.Model)
	if !ok {
		return "", errors.New("unexpected instance selection result")
	}
	if selectionModel.Err != nil {
		return "", selectionModel.Err
	}
	if selectionModel.Aborted {
		return "", errors.New("user aborted instance selection")
	}
	if selectionModel.Choice == "" {
		return "", errors.New("no instance selected")
	}

	return selectionModel.Choice, nil
}

func loadItems() selectionv2.Loader {
	return func() ([]listpkg.Item, error) {
		response, err := api.ListAllInstance()
		if err != nil {
			return nil, err
		}
		if len(response.Payload) == 0 {
			return nil, errors.New("no instances found")
		}

		items := make([]listpkg.Item, len(response.Payload))
		for i, instance := range response.Payload {
			items[i] = selectionv2.Item(instance.Name)
		}

		return items, nil
	}
}
