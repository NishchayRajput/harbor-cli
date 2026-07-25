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
package label

import (
	"errors"
	"fmt"

	listpkg "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/selectionv2"
)

func SelectLabelID(opts api.ListFlags) (int64, error) {
	labelIDs := make(map[string]int64)

	model := selectionv2.NewModel(
		"Label",
		"Loading labels...",
		loadLabelItems(opts, labelIDs),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return 0, fmt.Errorf("error running label selection: %w", err)
	}

	loadedModel, ok := finalModel.(selectionv2.Model)
	if !ok {
		return 0, errors.New("unexpected label selection model result")
	}

	if loadedModel.Err != nil {
		return 0, loadedModel.Err
	}

	labelID, ok := labelIDs[loadedModel.Choice]
	if !ok {
		return 0, fmt.Errorf("failed to get label id for %q", loadedModel.Choice)
	}

	return labelID, nil
}

func loadLabelItems(opts api.ListFlags, labelIDs map[string]int64) selectionv2.Loader {
	return func() ([]listpkg.Item, error) {
		response, err := api.ListLabel(opts)
		if err != nil {
			return nil, err
		}

		items := make([]listpkg.Item, len(response.Payload))
		for i, label := range response.Payload {
			labelIDs[label.Name] = label.ID
			items[i] = selectionv2.Item(label.Name)
		}

		return items, nil
	}
}
