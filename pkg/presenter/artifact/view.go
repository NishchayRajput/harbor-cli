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
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	artifactview "github.com/goharbor/harbor-cli/pkg/views/artifact/view"
)

type artifactLoadedMsg struct {
	artifact *models.Artifact
	err      error
}

type artifactLoaderModel struct {
	spinner  spinner.Model
	message  string
	loadFunc func() (*models.Artifact, error)
	artifact *models.Artifact
	err      error
}

func ViewArtifact(projectName, repoName, reference string) error {
	artifact, err := loadArtifact(
		fmt.Sprintf("Loading artifact details for %s/%s@%s...", projectName, repoName, shortReference(reference)),
		func() (*models.Artifact, error) {
			response, err := api.ViewArtifact(projectName, repoName, reference, false)
			if err != nil {
				return nil, err
			}

			return response.Payload, nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get info of an artifact: %w", err)
	}

	artifactview.ViewArtifact(artifact)
	return nil
}

func loadArtifact(message string, loadFunc func() (*models.Artifact, error)) (*models.Artifact, error) {
	spin := spinner.New()
	spin.Spinner = spinner.Dot

	model := artifactLoaderModel{
		spinner:  spin,
		message:  message,
		loadFunc: loadFunc,
	}

	finalModel, err := tea.NewProgram(model).Run()
	if err != nil {
		return nil, err
	}

	loadedModel, ok := finalModel.(artifactLoaderModel)
	if !ok {
		return nil, fmt.Errorf("unexpected artifact loader result")
	}

	return loadedModel.artifact, loadedModel.err
}

func (m artifactLoaderModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadCmd())
}

func (m artifactLoaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case artifactLoadedMsg:
		m.artifact = msg.artifact
		m.err = msg.err
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.err = fmt.Errorf("user aborted artifact loading")
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m artifactLoaderModel) View() string {
	return fmt.Sprintf("%s %s\n", m.spinner.View(), m.message)
}

func (m artifactLoaderModel) loadCmd() tea.Cmd {
	return func() tea.Msg {
		artifact, err := m.loadFunc()
		return artifactLoadedMsg{
			artifact: artifact,
			err:      err,
		}
	}
}

func shortReference(reference string) string {
	if len(reference) <= 24 {
		return reference
	}

	return strings.TrimSpace(reference[:24]) + "..."
}
