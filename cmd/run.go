package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"worker_pool/pool"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Структура для работы с терминальным интерфейсом
type model struct {
	cursor  int
	choice  string
	inputs  []textinput.Model
	frame   int
	spinner spinner.Model
	wp      *pool.WorkerPool
}

// Основа для стилизирования интерйфейса
var (
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	textStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("111"))
	noStyle       = lipgloss.NewStyle()
	focusedButton = cursorStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

// Варианты выбора функций в TUI
var choices = []string{"Добавить Jobs", "Добавить Workers", "Удалить Workers"}

// Запуск TUI
func Run(wp *pool.WorkerPool) {
	if _, err := tea.NewProgram(model{wp: wp}).Run(); err != nil {
		fmt.Println("О, нет:", err)
		os.Exit(1)
	}
}

// Инициализация модели для терминального интерфейса
func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
	)
}

// Обновление модели при взаимодействии с TUI
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.frame {
	case 0:
		return m.UpdateMainFrame(msg)
	case 1:
		return m.UpdateInputFrame(msg)
	}
	return m, nil
}

// Обновление главного фрейма
func (m model) UpdateMainFrame(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc": // Завершение программы
			return m, tea.Quit

		case "enter": // Запись выбранной строки в нашу модель.
			m.choice = choices[m.cursor]
			m.frame = 1
			m.cursor = 0
			m.getInputFrame()
			return m, nil
		case "down", "s": // Передвижение курсора вниз
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case "up", "w": // Передвижение курсора вверх
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

// Обновление фрейма ввода данных
func (m model) UpdateInputFrame(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit

		// Поставить курсор на следующее поле
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Если пользователь нажал enter когда кнопка отправки была на курсору
			if s == "enter" && m.cursor == len(m.inputs) {
				m.frame = 0
				m.ClickSubmitButton()
				return m, nil
			}

			// Движение вверх
			if s == "up" || s == "shift+tab" {
				m.cursor--
			} else {
				m.cursor++
			}

			// Движение вниз
			if m.cursor > len(m.inputs) {
				m.cursor = 0
			} else if m.cursor < 0 {
				m.cursor = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.cursor {
					// Установка курсора
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = cursorStyle
					m.inputs[i].TextStyle = cursorStyle
					continue
				}
				// Убираем курсор
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

// Запуск функций Пула при нажатии кнопки
func (m *model) ClickSubmitButton() {
	switch m.choice {
	case choices[0]: // AddJobs
		n, _ := strconv.Atoi(m.inputs[0].Value())
		s := m.inputs[1].Value()
		go m.wp.AddJobs(n, s)
	case choices[1]: // AddWorkers
		n, _ := strconv.Atoi(m.inputs[0].Value())
		go m.wp.AddWorkers(n)
	case choices[2]: // DeleteJobs
		n, _ := strconv.Atoi(m.inputs[0].Value())
		go m.wp.DeleteWorkers(n)
	}
}

// Получение выбранного фрейма
func (m *model) getInputFrame() {
	switch m.choice {
	case choices[0]: // AddJob
		m.getJobFrame()
	case choices[1]: // Worker
		m.getWorkerFrame()
	case choices[2]: // DeleteJob
		m.getWorkerFrame()
	}
}

// Получение фрейма джобов
func (m *model) getJobFrame() {
	m.inputs = make([]textinput.Model, 2)
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Введите количество job'ов"
			t.Focus()
			t.PromptStyle = cursorStyle
			t.TextStyle = cursorStyle
		case 1:
			t.Placeholder = "Введите информацию о job'ах"
			t.CharLimit = 64
		}

		m.inputs[i] = t
	}
}

// Получение фрейма воркеров
func (m *model) getWorkerFrame() {
	m.inputs = make([]textinput.Model, 1)
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32
		t.Placeholder = "Введите количество worker'ов"
		t.Focus()
		t.PromptStyle = cursorStyle
		t.TextStyle = cursorStyle
		m.inputs[i] = t
	}
}

// Полученние введенных строк в окне ввода
func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// Отрисовка TUI
func (m model) View() string {
	switch m.frame {
	case 0:
		return m.ViewMainFrame()
	case 1:
		return m.ViewInputFrame()
	}
	return ""
}

// Отрисовка главного фрейма
func (m model) ViewMainFrame() string {
	s := strings.Builder{}
	s.WriteString(textStyle.Render("\n"))
	s.WriteString(textStyle.Render("Дорогой Senior, Выберите функцию:)\n"))
	s.WriteString(textStyle.Render("\n"))
	for i := 0; i < len(choices); i++ {
		if m.cursor == i { // отрисовка текущего положения курсора
			s.WriteString(chosenView(choices[i]))
		} else {
			s.WriteString(choicesView(choices[i]))
		}
		s.WriteString("\n")
	}
	s.WriteString(m.subtleView())
	return mainStyle.Render(s.String())
}

// Отрисовка фрейма ввода
func (m model) ViewInputFrame() string {
	s := strings.Builder{}
	s.WriteString(textStyle.Render("\n"))
	s.WriteString(textStyle.Render("Дорогой Senior, Вы выбрали: "))
	s.WriteString(cursorStyle.Render("`" + m.choice + "`"))
	s.WriteString(textStyle.Render(" :)"))
	s.WriteString(textStyle.Render("\n\n"))
	for i := range m.inputs {
		s.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			s.WriteRune('\n')
		}
	}
	button := &blurredButton
	if m.cursor == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&s, "\n\n%s\n", *button)
	s.WriteString(m.subtleView())
	return mainStyle.Render(s.String())
}

// Отросиовка выбранной строки
func chosenView(label string) string {
	return cursorStyle.Render("[🍕] " + label)
}

// Отрисовка невыбранных строк
func choicesView(label string) string {
	return "[  ] " + label
}

// Отрисовка подсказок
func (m *model) subtleView() string {
	s := "\n"
	s += textStyle.Render("Текущее количество Worker:", strconv.Itoa(m.wp.GetWorkersCnt()))
	s += "\n"
	s += textStyle.Render("Текущее количество Job:", strconv.Itoa(m.wp.GetJobCnt()))
	s += "\n"
	s += subtleStyle.Render("w/s, up/down: движение курсора, ") +
		subtleStyle.Render("enter: выбор, ") +
		subtleStyle.Render("q, esc: выход")
	s += "\n"
	return s
}
