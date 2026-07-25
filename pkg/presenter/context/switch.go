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

	listpkg "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/selectionv2"
)

func SelectActiveContext() (string, error) {
	model := selectionv2.NewModel(
		"Context",
		"Loading contexts...",
		loadContextItems(),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return "", fmt.Errorf("error during context selection: %w", err)
	}

	selectionModel, ok := finalModel.(selectionv2.Model)
	if !ok {
		return "", errors.New("unexpected context selection result")
	}
	if selectionModel.Err != nil {
		return "", selectionModel.Err
	}
	if selectionModel.Aborted {
		return "", errors.New("user aborted context selection")
	}
	if selectionModel.Choice == "" {
		return "", errors.New("no context selected")
	}

	return selectionModel.Choice, nil
}

func loadContextItems() selectionv2.Loader {
	return func() ([]listpkg.Item, error) {
		config, err := utils.GetCurrentHarborConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get config: %w", err)
		}

		contexts := buildContextViews(config)
		if len(contexts) == 0 {
			return nil, errors.New("no contexts found")
		}

		items := make([]listpkg.Item, 0, len(contexts))
		for _, ctx := range contexts {
			name := "  " + ctx.Name
			if ctx.Name == config.CurrentCredentialName {
				name = "* " + ctx.Name
			}

			items = append(items, selectionv2.Item(name))
		}

		return items, nil
	}
}
