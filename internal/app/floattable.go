package app

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hedhyw/json-log-viewer/internal/pkg/source"
)

type floatTableModel struct {
	helper

	table          table.Model
	lastWindowSize tea.WindowSizeMsg

	logEntries source.LogEntries
}

func newFloatTableModel(application Application, logEntries source.LogEntries) floatTableModel {
	helper := helper{Application: application}

	const cellIDLogLevel = 1
	// var wsz tea.WindowSizeMsg = tea.WindowSizeMsg{
	// 	Width:  application.LastWindowSize.Width - 100,
	// 	Height: application.LastWindowSize.Height - 100,
	// }
	wsz := application.LastWindowSize
	tableLogs := table.New(
		table.WithColumns([]table.Column{{Width: wsz.Width, Title: "Sort By Field"}}),
		table.WithFocused(true),
		table.WithHeight(wsz.Height-100),
		table.WithWidth(wsz.Width-100),
	)

	//tableLogs.SetRows([]table.Row{logEntries[0].Row()})
	x := []table.Row{{"abc"}, {"def"}, {"arsalan"}}

	tableLogs.SetStyles(getTableStyles())
	tableLogs.SetRows(x)

	tableStyles := getTableStyles()
	tableStyles.RenderCell = func(_ table.Model, value string, position table.CellPosition) string {
		style := tableStyles.Cell

		if position.Column == cellIDLogLevel {
			return removeClearSequence(
				helper.getLogLevelStyle(
					logEntries,
					style,
					position.RowID,
				).Render(value),
			)
		}

		return style.Render(value)
	}

	tableLogs.SetStyles(tableStyles)

	return floatTableModel{
		helper:     helper,
		table:      tableLogs,
		logEntries: logEntries,
	}.handleWindowSizeMsg(wsz)
}

// Init initializes component. It implements tea.Model.
func (m floatTableModel) Init() tea.Cmd {
	return nil
}

// View renders component. It implements tea.Model.
func (m floatTableModel) View() string {

	return m.table.View()
}

// Update handles events. It implements tea.Model.
func (m floatTableModel) Update(msg tea.Msg) (floatTableModel, tea.Cmd) {
	var cmdBatch []tea.Cmd

	m.helper = m.helper.Update(msg)

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m = m.handleWindowSizeMsg(msg)
	}

	m.table, cmdBatch = batched(m.table.Update(msg))(cmdBatch)

	return m, tea.Batch(cmdBatch...)
}

func (m floatTableModel) handleWindowSizeMsg(msg tea.WindowSizeMsg) floatTableModel {
	const heightOffset = 4

	x, y := m.BaseStyle.GetFrameSize()
	m.table.SetWidth(msg.Width - x*2)
	m.table.SetHeight(msg.Height - y*2 - footerSize - heightOffset)
	//m.table.SetColumns(getColumns(m.table.Width()-10, m.Config))
	m.lastWindowSize = msg

	return m
}

// Cursor returns the index of the selected row.
func (m floatTableModel) Cursor() int {
	return m.table.Cursor()
}
