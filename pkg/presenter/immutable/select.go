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

	listpkg "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/selectionv2"
)

func SelectRule(projectName string) (int64, error) {
	lookup := map[string]int64{}
	model := selectionv2.NewModel(
		"Immutable Rule",
		fmt.Sprintf("Loading immutable rules for %s...", projectName),
		loadItems(projectName, lookup),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return 0, fmt.Errorf("error during immutable rule selection: %w", err)
	}

	selectionModel, ok := finalModel.(selectionv2.Model)
	if !ok {
		return 0, errors.New("unexpected immutable rule selection result")
	}
	if selectionModel.Err != nil {
		return 0, selectionModel.Err
	}
	if selectionModel.Aborted {
		return 0, errors.New("user aborted immutable rule selection")
	}
	if selectionModel.Choice == "" {
		return 0, errors.New("no immutable rule selected")
	}

	id, ok := lookup[selectionModel.Choice]
	if !ok {
		return 0, errors.New("selected immutable rule not found")
	}

	return id, nil
}

func loadItems(projectName string, lookup map[string]int64) selectionv2.Loader {
	return func() ([]listpkg.Item, error) {
		response, err := api.ListImmutable(projectName)
		if err != nil {
			return nil, err
		}
		if len(response.Payload) == 0 {
			return nil, errors.New("no immutable rules found")
		}

		items := make([]listpkg.Item, 0, len(response.Payload))
		for _, rule := range response.Payload {
			label := formatRuleLabel(rule)
			lookup[label] = rule.ID
			items = append(items, selectionv2.Item(label))
		}

		return items, nil
	}
}

func formatRuleLabel(rule *models.ImmutableRule) string {
	scopeSelectors := make([]string, 0, len(rule.ScopeSelectors))
	for _, scope := range rule.ScopeSelectors {
		for _, repo := range scope {
			scopeSelectors = append(scopeSelectors, fmt.Sprintf("%s %s", repo.Decoration, repo.Pattern))
		}
	}

	tagSelectors := make([]string, 0, len(rule.TagSelectors))
	for _, tag := range rule.TagSelectors {
		tagSelectors = append(tagSelectors, fmt.Sprintf("%s %s", tag.Decoration, tag.Pattern))
	}

	scope := "-"
	if len(scopeSelectors) > 0 {
		scope = scopeSelectors[0]
	}

	tag := "-"
	if len(tagSelectors) > 0 {
		tag = tagSelectors[0]
	}

	return fmt.Sprintf("for the %s, tags %s", scope, tag)
}
