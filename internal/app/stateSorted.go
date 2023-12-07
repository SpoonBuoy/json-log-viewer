package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

// StateSorted is a state that shows filtered records.
type StateSorted struct {
	helper

	previousState StateLoaded
	table         logsTableModel
	logEntries    source.LogEntries

	filterText string
}

func newStateSorted(
	application Application,
	previousState StateLoaded,
	filterText string,
) StateSorted {
	return StateSorted{
		helper: helper{Application: application},

		previousState: previousState,
		table:         previousState.table,

		filterText: filterText,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateSorted) Init() tea.Cmd {
	return func() tea.Msg {
		return events.LogEntriesLoadedMsg(
			s.previousState.logEntries.Sort(s.filterText),
		)
	}
}

// View renders component. It implements tea.Model.
func (s StateSorted) View() string {
	footer := s.Application.FooterStyle.Render(" sorted by: " + s.filterText)

	return s.BaseStyle.Render(s.table.View()) + "\n" + footer
}

// Update handles events. It implements tea.Model.
func (s StateSorted) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.helper = s.helper.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case events.BackKeyClickedMsg:
		return s.previousState.withApplication(s.Application)
	case events.FilterKeyClickedMsg:
		return s.handleFilterKeyClickedMsg()
	case events.EnterKeyClickedMsg, events.ArrowRightKeyClickedMsg:
		return s.handleRequestOpenJSON()
	case events.LogEntriesLoadedMsg:
		return s.handleLogEntriesLoadedMsg(msg)
	case events.OpenJSONRowRequestedMsg:
		return s.handleOpenJSONRowRequestedMsg(msg, s)
	case tea.KeyMsg:
		if cmd := s.handleKeyMsg(msg); cmd != nil {
			return s, cmd
		}
	default:
		s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)
	}

	s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)

	return s, tea.Batch(cmdBatch...)
}

func (s StateSorted) handleLogEntriesLoadedMsg(
	msg events.LogEntriesLoadedMsg,
) (tea.Model, tea.Cmd) {
	s.logEntries = source.LogEntries(msg)
	s.table = newLogsTableModel(s.Application, s.logEntries)

	return s, s.table.Init()
}

func (s StateSorted) handleFilterKeyClickedMsg() (tea.Model, tea.Cmd) {
	state := newStateFiltering(
		s.Application,
		s.previousState,
	)

	return initializeModel(state)
}

func (s StateSorted) handleRequestOpenJSON() (tea.Model, tea.Cmd) {
	if len(s.logEntries) == 0 {
		return s, events.BackKeyClicked
	}

	return s, events.OpenJSONRowRequested(s.logEntries, s.table.Cursor())
}

func (s StateSorted) withApplication(application Application) (state, tea.Cmd) {
	s.Application = application

	var cmd tea.Cmd
	s.table, cmd = s.table.Update(s.Application.LastWindowSize)

	return s, cmd
}

// String implements fmt.Stringer.
func (s StateSorted) String() string {
	return modelValue(s)
}
