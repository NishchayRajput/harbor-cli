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

package artifact

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelistv2"
)

var tagColumns = []table.Column{
	{Title: "Name", Width: tablelist.WidthL},
	{Title: "Pull Time", Width: tablelist.WidthL},
	{Title: "Push Time", Width: tablelist.WidthL},
}

func ListTags(projectName, repoName, reference string) error {
	model := tablelistv2.NewModel(
		tagColumns,
		fmt.Sprintf("Loading tags for %s/%s@%s...", projectName, repoName, reference),
		loadTagRows(projectName, repoName, reference),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running tag list: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected tag list model result")
	}

	return loadedModel.Err
}

func loadTagRows(projectName, repoName, reference string) tablelistv2.Loader {
	return func() ([]table.Row, error) {
		response, err := api.ListTags(projectName, repoName, reference)
		if err != nil {
			return nil, fmt.Errorf("failed to list tags: %v", err)
		}

		return buildTagRows(response.Payload), nil
	}
}

func buildTagRows(tags []*models.Tag) []table.Row {
	rows := make([]table.Row, 0, len(tags))

	for _, tag := range tags {
		pullTime, _ := utils.FormatCreatedTime(tag.PullTime.String())
		pushTime, _ := utils.FormatCreatedTime(tag.PushTime.String())
		rows = append(rows, table.Row{
			tag.Name,
			pullTime,
			pushTime,
		})
	}

	return rows
}
