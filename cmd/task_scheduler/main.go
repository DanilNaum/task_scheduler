package main

import (
	"DanilNaum/task_scheduler/internal/testtasks"
	"bufio"
	"time"
	"os"
	"strings"
	"sync"
)

func main() {
	var second_task []string
	var minute_task []string
	var hour_task []string
	mu := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	numberOfSimultaneousRequests := 5
	data_channel := make(chan string)
	go func(){
		defer close(data_channel)
		for{
			dir, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			file, err := os.Open(dir + "/cmd/task_scheduler/tasks.txt")
			if err != nil {
				panic(err)
			}
			
			s := bufio.NewScanner(file)
			var second_task_tmp []string
			var minute_task_tmp []string
			var hour_task_tmp []string
			for s.Scan(){
				if strings.Split(s.Text()," ")[0] == "secondly"{
					second_task_tmp = append(second_task_tmp, s.Text())
				}else if strings.Split(s.Text()," ")[0] == "minutely"{
					minute_task_tmp = append(minute_task_tmp, s.Text())
				}else if strings.Split(s.Text()," ")[0] == "hourly"{
					hour_task_tmp = append(hour_task_tmp, s.Text())
				}
			}
			mu.Lock()
			second_task = second_task_tmp
			minute_task = minute_task_tmp
			hour_task = hour_task_tmp
			mu.Unlock()
			file.Close()
		}
	}()
	go func() {
		
		tickSecond := time.NewTicker(time.Second)
		tickMinute := time.NewTicker(time.Minute)
		tickHour := time.NewTicker(time.Hour)
		defer tickSecond.Stop() // освободим ресурсы, при завершении работы функции
		defer tickMinute.Stop()
		defer tickHour.Stop()
		for {
			select {
				case <-tickSecond.C:
					mu.Lock()
					for _,task :=range (second_task){
						data_channel<-task
					}
					mu.Unlock()
				case <-tickMinute.C:
					mu.Lock()
					for _,task :=range (minute_task){
						data_channel<-task
					}
					mu.Unlock()
				case <-tickHour.C:
					mu.Lock()
					for _,task :=range (hour_task){
						data_channel<-task
					}
					mu.Unlock()
			}
		}
	}()

	channel := make(chan []string)
	channel2 := make(chan struct{}, numberOfSimultaneousRequests)
	
	go func(){
		defer close(channel)
		defer close(channel2)
		for s := range data_channel{
			a := strings.Split(s," ")
			channel2 <- struct{}{} // канал для проверки завершения выполнения
				//  предыдущей задачи и ограничения количества одновреммено выполняемых задач
			channel <- a	// канал для передачи следующей задачи в фукцию
		}
	}()
	
	// Наличие структуры на thread означает, что он свободен и можно передать на него задачу
	thread := make(map[int] chan struct{})
	for i := 0; i < 5; i++ {
		thread[i] = make(chan struct{},1)
		thread[i] <-  struct{}{}
		defer close(thread[i])
	}
		
	for a := range channel{
		wg.Add(1)
		go func(wg *sync.WaitGroup, a []string){
			defer wg.Done()
			select {
			case <-thread[0]:
				testtasks.Wait(1,a[2:])
				thread[0] <-  struct{}{}
			case <-thread[1]:
				testtasks.Wait(2,a[2:])
				thread[1] <-  struct{}{}	
			case <-thread[2]:
				testtasks.Wait(3,a[2:])
				thread[2] <-  struct{}{}
			case <-thread[3]:
				testtasks.Wait(4,a[2:])
				thread[3] <-  struct{}{}	
			case <-thread[4]:
				testtasks.Wait(5,a[2:])
				thread[4] <-  struct{}{}	
			}
			<- channel2
			
		}(wg,a)
	}
	wg.Wait()
}
