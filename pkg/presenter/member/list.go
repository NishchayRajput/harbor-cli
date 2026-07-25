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

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelistv2"
	"github.com/sahilm/fuzzy"
)

var wideColumns = []table.Column{
	{Title: "ID", Width: 4},
	{Title: "Member Name", Width: 12},
	{Title: "Type", Width: 8},
	{Title: "Role Name", Width: 16},
	{Title: "Role ID", Width: 8},
	{Title: "Project ID", Width: 12},
}

var compactColumns = []table.Column{
	{Title: "ID", Width: 4},
	{Title: "Member Name", Width: 12},
	{Title: "Type", Width: 8},
	{Title: "Role Name", Width: 16},
}

func List(opts api.ListMemberOptions, searchQuery string, wide bool) error {
	columns := compactColumns
	if wide {
		columns = wideColumns
	}

	model := tablelistv2.NewModel(
		columns,
		fmt.Sprintf("Loading members for %s...", opts.ProjectNameOrID),
		loadRows(opts, searchQuery, wide),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running member list: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected member list model result")
	}

	return loadedModel.Err
}

func loadRows(opts api.ListMemberOptions, searchQuery string, wide bool) tablelistv2.Loader {
	return func() ([]table.Row, error) {
		response, err := api.ListMember(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get members list: %v", err)
		}
		if searchQuery != "" && opts.EntityName == "" {
			response.Payload = filterMembers(response.Payload, searchQuery)
		}
		if len(response.Payload) == 0 {
			return nil, errors.New("no members found")
		}

		return buildRows(response.Payload, wide), nil
	}
}

func filterMembers(members []*models.ProjectMemberEntity, query string) []*models.ProjectMemberEntity {
	set := make([]string, 0, len(members))
	for _, member := range members {
		set = append(set, member.EntityName)
	}

	matches := fuzzy.Find(query, set)
	results := make([]*models.ProjectMemberEntity, 0, len(matches))
	for _, match := range matches {
		results = append(results, members[match.Index])
	}

	return results
}

func buildRows(members []*models.ProjectMemberEntity, wide bool) []table.Row {
	rows := make([]table.Row, 0, len(members))

	for _, member := range members {
		memberID := strconv.FormatInt(member.ID, 10)
		roleName := utils.CamelCaseToHR(member.RoleName)
		memberType := member.EntityType
		if memberType == "u" {
			memberType = "User"
		} else if memberType == "g" {
			memberType = "Group"
		}

		if wide {
			rows = append(rows, table.Row{
				memberID,
				member.EntityName,
				memberType,
				roleName,
				strconv.FormatInt(member.RoleID, 10),
				strconv.FormatInt(member.ProjectID, 10),
			})
			continue
		}

		rows = append(rows, table.Row{
			memberID,
			member.EntityName,
			memberType,
			roleName,
		})
	}

	return rows
}
