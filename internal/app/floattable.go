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

	tableLogs := table.New(
		table.WithColumns([]table.Column{{Width: application.LastWindowSize.Width, Title: "Sort By Field"}}),
		table.WithFocused(true),
		table.WithHeight(application.LastWindowSize.Height-100),
		table.WithWidth(application.LastWindowSize.Width-100),
	)

	//read sort field options from config.Fields
	cfg := application.Config
	var rows []table.Row
	for _, field := range cfg.Fields {
		rows = append(rows, table.Row{field.Title})
	}

	tableLogs.SetStyles(getTableStyles())
	tableLogs.SetRows(rows)

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
	}.handleWindowSizeMsg(application.LastWindowSize)
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
