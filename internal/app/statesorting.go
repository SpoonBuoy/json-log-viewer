package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
)

// StateSorting is a state to prompt for filter term.
type StateSorting struct {
	helper
	initCmd       tea.Cmd
	previousState StateLoaded
	sortTable     floatTableModel
	sortByField   string
}

func newStateSorting(
	application Application,
	previousState StateLoaded,
) StateSorting {

	st := newFloatTableModel(application, previousState.logEntries)
	return StateSorting{
		helper:        helper{Application: application},
		initCmd:       st.Init(),
		previousState: previousState,
		//table:         previousState.table,
		sortTable: st,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateSorting) Init() tea.Cmd {
	return s.initCmd
}

// View renders component. It implements tea.Model.
func (s StateSorting) View() string {
	footer := s.Application.FooterStyle.Render("[A] Ascending; [Enter] Descending")
	return s.BaseStyle.Render(s.sortTable.View()) + "\n" + s.FooterStyle.Render(footer)
}

// Update handles events. It implements tea.Model.
func (s StateSorting) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.helper = s.helper.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case events.BackKeyClickedMsg:
		return s.previousState.withApplication(s.Application)
	case events.EnterKeyClickedMsg:
		return s.handleEnterKeyClickedMsg()
	case events.RevSortKeyClickedMsg:
		return s.handleRevSortKeyClickedMsg()
	case tea.KeyMsg:
		cmdBatch = append(cmdBatch, s.handleKeyMsg(msg)...)

		if s.isFilterKeyMap(msg) {
			// Intercept table update.
			return s, tea.Batch(cmdBatch...)
		}
	}

	s.sortTable, cmdBatch = batched(s.sortTable.Update(msg))(cmdBatch)

	return s, tea.Batch(cmdBatch...)
}

func (s StateSorting) handleKeyMsg(msg tea.KeyMsg) []tea.Cmd {
	var cmdBatch []tea.Cmd

	cmdBatch = appendCmd(cmdBatch, s.helper.handleKeyMsg(msg))

	if s.isArrowUpKeyMap(msg) {
		cmdBatch = appendCmd(cmdBatch, s.handleArrowUpKeyClicked())
	}

	return cmdBatch
}

func (s StateSorting) handleArrowUpKeyClicked() tea.Cmd {
	if s.sortTable.Cursor() == 0 {
		return events.ViewRowsReloadRequested
	}

	return nil
}
func (s StateSorting) handleEnterKeyClickedMsg() (tea.Model, tea.Cmd) {

	s.sortByField = getFieldFromConfigByIndex(s.sortTable.Cursor(), s.Config)
	return initializeModel(newStateSorted(
		s.Application,
		s.previousState,
		s.sortByField,
		true,
	))
}

func (s StateSorting) handleRevSortKeyClickedMsg() (tea.Model, tea.Cmd) {

	s.sortByField = getFieldFromConfigByIndex(s.sortTable.Cursor(), s.Config)
	return initializeModel(newStateSorted(
		s.Application,
		s.previousState,
		s.sortByField,
		false,
	))
}

// String implements fmt.Stringer.
func (s StateSorting) String() string {
	return modelValue(s)
}
