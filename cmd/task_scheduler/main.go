package main

import (
	"DanilNaum/task_scheduler/internal/testtasks"
	"bufio"
	"fmt"
	"os"
	"sync"
	"strings"
)



func main() {
	dir, err := os.Getwd()
    if err != nil {
        panic(err)
    }
	fmt.Print(dir)
	file, err := os.Open(dir + "/cmd/task_scheduler/tasks.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	wg := new(sync.WaitGroup)
	numberOfSimultaneousRequests := 2
	channel := make(chan string)
	channel2 := make(chan struct{}, numberOfSimultaneousRequests)
	
	go func(){
		defer close(channel)
		defer close(channel2)
		for s.Scan(){
			channel2 <- struct{}{} // канал для проверки завершения выполнения
			//  предыдущей задачи и ограничения количества одновреммено выполняемых задач
			channel <- s.Text()	// канал для передачи следующей задачи в фукцию
		}

		// close(channel)
		// close(channel2)
	}()
	for t := range channel{
		// t :=s.Text()
		a := strings.Split(t," ")
		// fmt.Println(a)
		// a := []string{"abc", "1", "5"}
		wg.Add(1)
		go func(wg *sync.WaitGroup){
			defer wg.Done()
			testtasks.Wait(a[2:])
			<- channel2
		}(wg)
	}
	wg.Wait()
	

}
