package cmd

import (
	"fmt"
	"worker_pool/pool"
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
		fmt.Println(" `1 <count> <data>` либо `add_job <count> <data>` - создаст новые job'ы с информацией `data` в количестве `count`,")
		fmt.Println(" `2 <count>` либо `add_worker <count>` - cоздаст новых worker'ов в количестве `count`,")
		fmt.Println(" `3 <count>` либо `delete_worker <count>` - удалит worker'ов в количестве `count`,")
		fmt.Println(" `4` либо `stop` - принудительно завершить программу")
	}

	fmt.Println(ColorGreen, "Привет, сеньор! Спасибо, что запустил меня ;).")
	showGoodInput()
	fmt.Print(ColorBlue)

	// Интерактивный режим ввода
	fmt.Scan(&s)
	for ; ; fmt.Scan(&s) {
		switch s {
		case "1", "add_job":
			fmt.Scan(&n, &s)
			fmt.Println(ColorGreen, "- Добавление", n, "джобов со строкой", s, "...")
			go wp.AddJobs(n, s)
		case "2", "add_worker":
			fmt.Scan(&n)
			fmt.Println(ColorGreen, "- Добавление", n, "воркеров...")
			go wp.AddWorkers(n)
		case "3", "delete_worker":
			fmt.Scan(&n)
			fmt.Println(ColorGreen, "- Удаление", n, "воркеров...")
			go wp.DeleteWorkers(n)
		case "4", "stop":
			return
		default:
			fmt.Println(ColorRed, "Неверный формат ввода. Не расстраивайся", ColorReset)
			showGoodInput()
		}
		fmt.Print(ColorBlue)
	}
}
