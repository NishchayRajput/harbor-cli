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

package member

import (
	"errors"
	"fmt"
	"strconv"

	listpkg "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/selectionv2"
)

func SelectRole() (int64, error) {
	roles := []string{"Project Admin", "Developer", "Guest", "Maintainer", "Limited Guest"}
	model := selectionv2.NewModel(
		"Role",
		"Loading roles...",
		loadRoleItems(roles),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return 0, fmt.Errorf("error during role selection: %w", err)
	}

	selectionModel, ok := finalModel.(selectionv2.Model)
	if !ok {
		return 0, errors.New("unexpected role selection result")
	}
	if selectionModel.Err != nil {
		return 0, selectionModel.Err
	}
	if selectionModel.Aborted {
		return 0, errors.New("user aborted role selection")
	}
	if selectionModel.Choice == "" {
		return 0, errors.New("no role selected")
	}

	for i, role := range roles {
		if role == selectionModel.Choice {
			return int64(i + 1), nil
		}
	}

	return 0, errors.New("selected role not found")
}

func SelectMember(projectName, memberName string) (int64, error) {
	lookup := map[string]int64{}
	model := selectionv2.NewModel(
		"Member",
		fmt.Sprintf("Loading members for %s...", projectName),
		loadMemberItems(projectName, memberName, lookup),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return 0, fmt.Errorf("error during member selection: %w", err)
	}

	selectionModel, ok := finalModel.(selectionv2.Model)
	if !ok {
		return 0, errors.New("unexpected member selection result")
	}
	if selectionModel.Err != nil {
		return 0, selectionModel.Err
	}
	if selectionModel.Aborted {
		return 0, errors.New("user aborted member selection")
	}
	if selectionModel.Choice == "" {
		return 0, errors.New("no member selected")
	}

	id, ok := lookup[selectionModel.Choice]
	if !ok {
		return 0, errors.New("selected member not found")
	}

	return id, nil
}

func loadRoleItems(roles []string) selectionv2.Loader {
	return func() ([]listpkg.Item, error) {
		items := make([]listpkg.Item, len(roles))
		for i, role := range roles {
			items[i] = selectionv2.Item(role)
		}
		return items, nil
	}
}

func loadMemberItems(projectName, memberName string, lookup map[string]int64) selectionv2.Loader {
	return func() ([]listpkg.Item, error) {
		response, err := api.ListMembers(projectName, memberName, true)
		if err != nil {
			return nil, err
		}
		if len(response.Payload) == 0 {
			return nil, errors.New("no members found")
		}

		items := make([]listpkg.Item, 0, len(response.Payload))
		for _, member := range response.Payload {
			label := member.EntityName + " (" + strconv.FormatInt(member.ID, 10) + ")"
			lookup[label] = member.ID
			items = append(items, selectionv2.Item(label))
		}

		return items, nil
	}
}
