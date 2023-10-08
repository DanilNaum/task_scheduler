package testtasks

import (
	"fmt"
	"math/rand"
	"strconv"
	// "sync"
	"time"
)
func Wait(args []string){
	if len(args) == 3 && args[1]<args[2]{
		
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "thread ",": starting job",args[0])
		min_time,err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
		max_time,err := strconv.Atoi(args[2])
		if err != nil {
			panic(err)
		}
		dalay_time := rand.Intn( max_time-min_time+1) + min_time
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "thread ",": wating for ", dalay_time)
		time.Sleep(time.Duration(dalay_time)*time.Second)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), "thread ",":job is complited")
		

	} else {
		fmt.Println("Неверные аргументы функции")
// Обработка передачи неверных аргументов функции, выбор способа реакции зависит от цели
		// panic("Неверные аргументы функции")
	}
}