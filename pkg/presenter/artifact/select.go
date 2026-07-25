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

	listpkg "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/base/selectionv2"
)

func SelectReference(projectName, repoName string) (string, error) {
	model := selectionv2.NewModel(
		"Artifact",
		fmt.Sprintf("Loading artifacts for %s/%s...", projectName, repoName),
		loadReferenceItems(projectName, repoName),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return "", fmt.Errorf("error during artifact selection: %w", err)
	}

	selectionModel, ok := finalModel.(selectionv2.Model)
	if !ok {
		return "", errors.New("unexpected artifact selection result")
	}
	if selectionModel.Err != nil {
		return "", selectionModel.Err
	}
	if selectionModel.Aborted {
		return "", errors.New("user aborted artifact selection")
	}
	if selectionModel.Choice == "" {
		return "", errors.New("no artifact selected")
	}

	return selectionModel.Choice, nil
}

func SelectTag(projectName, repoName, reference string) (string, error) {
	model := selectionv2.NewModel(
		"Tag",
		fmt.Sprintf("Loading tags for %s/%s@%s...", projectName, repoName, reference),
		loadTagItems(projectName, repoName, reference),
	)

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return "", fmt.Errorf("error during tag selection: %w", err)
	}

	selectionModel, ok := finalModel.(selectionv2.Model)
	if !ok {
		return "", errors.New("unexpected tag selection result")
	}
	if selectionModel.Err != nil {
		return "", selectionModel.Err
	}
	if selectionModel.Aborted {
		return "", errors.New("user aborted tag selection")
	}
	if selectionModel.Choice == "" {
		return "", errors.New("no tag selected")
	}

	return selectionModel.Choice, nil
}

func loadReferenceItems(projectName, repoName string) selectionv2.Loader {
	return func() ([]listpkg.Item, error) {
		response, err := api.ListArtifact(projectName, repoName)
		if err != nil {
			return nil, err
		}
		if len(response.Payload) == 0 {
			return nil, errors.New("no artifacts found")
		}

		items := make([]listpkg.Item, len(response.Payload))
		for i, artifact := range response.Payload {
			items[i] = selectionv2.Item(artifact.Digest)
		}

		return items, nil
	}
}

func loadTagItems(projectName, repoName, reference string) selectionv2.Loader {
	return func() ([]listpkg.Item, error) {
		response, err := api.ListTags(projectName, repoName, reference)
		if err != nil {
			return nil, err
		}
		if len(response.Payload) == 0 {
			return nil, errors.New("no tags found")
		}

		items := make([]listpkg.Item, len(response.Payload))
		for i, tag := range response.Payload {
			items[i] = selectionv2.Item(tag.Name)
		}

		return items, nil
	}
}
