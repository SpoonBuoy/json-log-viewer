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
	//table         logsTableModel
	sortTable floatTableModel
	//textInput textinput.Model
}

func newStateSorting(
	application Application,
	previousState StateLoaded,
) StateSorting {
	//textInput := textinput.New()
	//textInput.Focus()
	st := newFloatTableModel(application, previousState.logEntries)
	return StateSorting{
		helper:        helper{Application: application},
		initCmd:       st.Init(),
		previousState: previousState,
		//table:         previousState.table,
		sortTable: st,
		//textInput: textInput,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateSorting) Init() tea.Cmd {
	return s.initCmd
}

// View renders component. It implements tea.Model.
func (s StateSorting) View() string {
	//t := newFloatTableModel(s.Application, s.sortTable.logEntries)
	//return s.BaseStyle.Render(s.sortTable.View()) + s.sortTable.View() + "\n" + s.textInput.View()
	return s.BaseStyle.Render(s.sortTable.View())
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
	// case tea.KeyMsg:
	// 	if cmd := s.handleKeyMsg(msg); cmd != nil {
	// 		// Intercept table update.
	// 		return s, cmd
	// 	}

	default:
		s.sortTable, cmdBatch = batched(s.sortTable.Update(msg))(cmdBatch)
	}

	//s.textInput, cmdBatch = batched(s.textInput.Update(msg))(cmdBatch)

	return s, tea.Batch(cmdBatch...)
}

// func (s StateSorting) handleKeyMsg(msg tea.KeyMsg) tea.Cmd {
// 	if len(msg.Runes) == 1 {
// 		return nil
// 	}

// 	return s.helper.handleKeyMsg(msg)
// }

func (s StateSorting) handleEnterKeyClickedMsg() (tea.Model, tea.Cmd) {
	// if s.textInput.Value() == "" {
	// 	return s, events.BackKeyClicked
	// }

	return initializeModel(newStateSorted(
		s.Application,
		s.previousState,
		"xxx",
		//s.textInput.Value(),
	))
}

// String implements fmt.Stringer.
func (s StateSorting) String() string {
	return modelValue(s)
}
