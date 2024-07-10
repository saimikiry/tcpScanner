package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"time"
)

func startWorker(host string, ports_pool <-chan int, open_ports chan<- int) {
	for port := range ports_pool {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.Dial("tcp", address)
		if err == nil {
			conn.Close()
			open_ports <- port
		} else {
			open_ports <- 0
		}
	}
}

func main() {
	// Получение аргументов
	args := os.Args[1:]

	// Проверка числа аргументов
	if len(args) != 5 {
		fmt.Printf("[Scanner] Incorrect import! Usage example:\n[Scanner] $ scanner.exe <host> <first port> <last port> <goroutins count>\n[Scanner] $ scanner.exe scanme.nmap.org 1 1024 100")
		return
	}

	// Фиксация хоста
	host := string(args[0])

	// Фиксация диапазона портов
	first_port, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("[Scanner] Incorrent first port!")
		return
	}

	last_port, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("[Scanner] Incorrect last port!")
		return
	}

	// Фиксация числа горутин
	workers_count, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Printf("[Scanner] Incorrect workers count!")
		return
	}

	duration, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Printf("[Scanner] Incorrect time duration!")
		return
	}

	// Создание буферизованного канала для передачи портов воркерам
	ports_pool := make(chan int, workers_count)
	defer close(ports_pool)

	// Создание небуферизованного канала для получения результатов сканирования
	open_ports := make(chan int)
	defer close(open_ports)

	// Запуск горутин для проведения сканирования
	for i := 0; i < workers_count; i++ {
		go startWorker(host, ports_pool, open_ports)
	}

	// Подача портов для сканирования
	go func() {
		for i := first_port; i <= last_port; i++ {
			ports_pool <- i
			time.Sleep(time.Duration(duration) * time.Millisecond)
			fmt.Printf("Pushed %d\n", i)
		}
	}()

	// Ожидание результатов сканирования
	result := []int{}
	for i := first_port; i <= last_port; i++ {
		cur := <-open_ports
		if cur != 0 {
			result = append(result, cur)
		}
	}

	// Сортировка и вывод результатов
	sort.Ints(result)
	fmt.Print(result)
}
