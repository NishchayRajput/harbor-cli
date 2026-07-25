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

package immutable

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelistv2"
)

var columns = []table.Column{
	{Title: "ID", Width: 12},
	{Title: "Repository", Width: 30},
	{Title: "Tag", Width: 30},
}

func ListRules(projectName string) error {
	model := tablelistv2.NewModel(
		columns,
		fmt.Sprintf("Loading immutable rules for %s...", projectName),
		loadRows(projectName),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running immutable rule list: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected immutable rule list model result")
	}

	return loadedModel.Err
}

func loadRows(projectName string) tablelistv2.Loader {
	return func() ([]table.Row, error) {
		response, err := api.ListImmutable(projectName)
		if err != nil {
			return nil, fmt.Errorf("failed to list immutablility rule: %v", err)
		}
		if len(response.Payload) == 0 {
			return nil, errors.New("no immutable tag rules found")
		}

		return buildRows(response.Payload), nil
	}
}

func buildRows(rules []*models.ImmutableRule) []table.Row {
	rows := make([]table.Row, 0, len(rules))

	for _, rule := range rules {
		var scopeSelectors []string
		for _, scope := range rule.ScopeSelectors {
			for _, repo := range scope {
				scopeSelectors = append(scopeSelectors, fmt.Sprintf("%s %s", repo.Decoration, repo.Pattern))
			}
		}

		tagSelectors := make([]string, len(rule.TagSelectors))
		for i, tag := range rule.TagSelectors {
			tagSelectors[i] = fmt.Sprintf("%s %s", tag.Decoration, tag.Pattern)
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", rule.ID),
			strings.Join(scopeSelectors, " "),
			strings.Join(tagSelectors, " "),
		})
	}

	return rows
}
