package filepickermodel

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type FilepickerModelConfig struct {
	Title            string
	CurrentDirectory string
	AllowedTypes     []string
	EnableFastSelect bool
	InfoBubble       string
}

type FilepickerModel struct {
	filepicker       filepicker.Model
	quitting         bool
	err              error
	enableFastSelect bool
	title            string
	SelectedFile     string
	infoBubble       string
	filepickerWidth  int
	width            int
}

type clearErrorMsg struct{}

var (
	docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m FilepickerModel) Width(width int) FilepickerModel {
	m.width = width
	return m
}

func (m FilepickerModel) Height(v int) FilepickerModel {
	m.filepicker.Height = v
	m.filepicker.AutoHeight = false
	return m
}

func (m FilepickerModel) FilepickerWidth(v int) FilepickerModel {
	m.filepickerWidth = v
	return m
}

func (m FilepickerModel) AutoHeight(v bool) FilepickerModel {
	m.filepicker.AutoHeight = v
	return m
}

func (m FilepickerModel) ShowSize(v bool) FilepickerModel {
	m.filepicker.ShowSize = v
	return m
}

func (m FilepickerModel) ShowHidden(v bool) FilepickerModel {
	m.filepicker.ShowHidden = v
	return m
}

func (m FilepickerModel) ShowPermissions(v bool) FilepickerModel {
	m.filepicker.ShowPermissions = v
	return m
}

func (m FilepickerModel) StyleDisabledCursor(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.DisabledCursor = v
	return m
}

func (m FilepickerModel) StyleCursor(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.Cursor = v
	return m
}

func (m FilepickerModel) StyleSymlink(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.Symlink = v
	return m
}

func (m FilepickerModel) StyleDirectory(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.Directory = v
	return m
}

func (m FilepickerModel) StyleFile(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.File = v
	return m
}

func (m FilepickerModel) StyleDisabledFile(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.DisabledFile = v
	return m
}

func (m FilepickerModel) StylePermission(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.Permission = v
	return m
}

func (m FilepickerModel) StyleSelected(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.Selected = v
	return m
}

func (m FilepickerModel) StyleDisabledSelected(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.DisabledSelected = v
	return m
}

func (m FilepickerModel) StyleFileSize(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.FileSize = v
	return m
}

func (m FilepickerModel) StyleEmptyDirectory(v lipgloss.Style) FilepickerModel {
	m.filepicker.Styles.EmptyDirectory = v
	return m
}

func (m FilepickerModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m FilepickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.SelectedFile = path
		if m.enableFastSelect {
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.SelectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m FilepickerModel) View() string {
	if m.quitting {
		return ""
	}

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.SelectedFile == "" {
		title := "Pick a file:"
		if len(m.title) > 0 {
			title = m.title
		}
		s.WriteString(title)
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.SelectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")

	var doc strings.Builder

	// Form
	{
		var formView string
		var infoView string

		// Form / Form
		{
			var s strings.Builder
			s.WriteString("\n  ")
			if m.err != nil {
				s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
			} else if m.SelectedFile == "" {
				title := "Pick a file:"
				if len(m.title) > 0 {
					title = m.title
				}
				s.WriteString(title)
			} else {
				s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.SelectedFile))
			}
			s.WriteString("\n\n" + m.filepicker.View() + "\n")
			if m.filepickerWidth > 0 {
				formView = lipgloss.NewStyle().Width(m.filepickerWidth).Render(s.String())
			} else {
				formView = s.String()
			}
		}

		// Form / Info
		{
			if m.infoBubble != "" && (m.width == 0 || (m.width > 0 && physicalWidth >= m.width*4/5)) {
				infoView = m.infoBubble
			}
		}

		doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, formView, infoView))
		doc.WriteString("\n\n")
	}

	if physicalWidth > 0 {
		docStyle = docStyle.MaxWidth(physicalWidth)
	}

	return docStyle.Render(doc.String())
}

func NewFilepickerModel(config FilepickerModelConfig) FilepickerModel {
	var dir string
	if len(config.CurrentDirectory) > 0 {
		dir = config.CurrentDirectory
	} else {
		dir, _ = os.Getwd()
	}

	fp := filepicker.New()
	fp.AllowedTypes = config.AllowedTypes
	fp.CurrentDirectory = dir

	m := FilepickerModel{
		filepicker:       fp,
		enableFastSelect: config.EnableFastSelect,
		title:            config.Title,
		infoBubble:       config.InfoBubble,
	}

	return m
}

func NewFilepickerModelProgram(config FilepickerModelConfig) *tea.Program {
	m := NewFilepickerModel(config)

	return tea.NewProgram(&m)
}
