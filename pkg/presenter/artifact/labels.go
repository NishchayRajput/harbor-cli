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

var labelColumns = []table.Column{
	{Title: "ID", Width: tablelist.WidthXS},
	{Title: "Name", Width: tablelist.WidthM},
	{Title: "Color", Width: tablelist.WidthM},
	{Title: "Description", Width: tablelist.WidthXL},
	{Title: "Creation Time", Width: tablelist.WidthL},
}

func ListLabels(projectName, repoName, reference string) error {
	model := tablelistv2.NewModel(
		labelColumns,
		fmt.Sprintf("Loading labels for %s/%s@%s...", projectName, repoName, reference),
		loadLabelRows(projectName, repoName, reference),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return fmt.Errorf("error running label list: %w", err)
	}

	loadedModel, ok := finalModel.(tablelistv2.Model)
	if !ok {
		return errors.New("unexpected label list model result")
	}

	return loadedModel.Err
}

func loadLabelRows(projectName, repoName, reference string) tablelistv2.Loader {
	return func() ([]table.Row, error) {
		response, err := api.ViewArtifact(projectName, repoName, reference, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get info of an artifact: %v", utils.ParseHarborErrorMsg(err))
		}

		if len(response.Payload.Labels) == 0 {
			return nil, fmt.Errorf("no labels found for artifact %s/%s@%s", projectName, repoName, reference)
		}

		return buildLabelRows(response.Payload.Labels), nil
	}
}

func buildLabelRows(labels []*models.Label) []table.Row {
	rows := make([]table.Row, 0, len(labels))

	for _, label := range labels {
		createdTime, _ := utils.FormatCreatedTime(label.CreationTime.String())
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", label.ID),
			label.Name,
			label.Color,
			label.Description,
			createdTime,
		})
	}

	return rows
}
