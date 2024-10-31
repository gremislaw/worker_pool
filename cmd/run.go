package cmd

import (
	"fmt"
	"worker_pool/pool"
)

// Нельзя менять. Нужно для понимания какой запрос пришел
const (
	add_job = "1"
	add_worker = "2"
	delete_worker = "3"
	stop = "4"
)

// Для разноцветного терминала

const (
	ColorReset  = "\033[0m"  // Сброс цвета
	ColorRed    = "\033[31m" // Красный
	ColorGreen  = "\033[32m" // Зеленый
	ColorYellow = "\033[33m" // Желтый
	ColorBlue   = "\033[34m" // Синий
)

func Run(wp *pool.WorkerPool) {
	s, n := "", 0

	// Выводит информацию о допустомом формате ввода
	showGoodInput := func() {
		fmt.Println(ColorYellow, "Допустимый формат ввода:")
		fmt.Println(ColorBlue, "1 <data> - добавить job с информацией data,")
		fmt.Println(" 2 <id> - добавить worker с номером id,")
		fmt.Println(" 3 <id> - удалить worker с номером id")
		fmt.Println(" 4 - принудительно завершить программу", ColorReset)
	}

	fmt.Println(ColorGreen, "Привет, сеньор! Спасибо, что запустил меня ;).")
	showGoodInput()

	// Интерактивный режим ввода
	fmt.Scan(&s)
	for ; ; fmt.Scan(&s) {
		switch s {
		case add_job:
			fmt.Scan(&s)
			go wp.AddJob(s)
		case add_worker:
			fmt.Scan(&n)
			go wp.AddWorker(n)
		case delete_worker:
			fmt.Scan(&n)
			go wp.DeleteWorker(n)
		case stop:
			return
		default:
			fmt.Println(ColorRed, "Неверный формат ввода. Не расстраивайся", ColorReset)
			showGoodInput()
		}
	}
}
