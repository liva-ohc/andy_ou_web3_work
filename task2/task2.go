package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"text/tabwriter"
	"time"
)

func main() {
	myInt := 10
	pointerTest1(&myInt)
	fmt.Println(myInt)

	mySlice := []int{1, 2, 3, 4, 5}
	pointerTest2(&mySlice)
	fmt.Println(mySlice)
	GoroutineTest1()

	fmt.Println("===================================================")
	GoroutineTest2()
	fmt.Println("===================================================")

	myRectangle := Rectangle{Width: 10, Height: 20}
	myCircle := Circle{Radius: 5}
	fmt.Println("矩形的面积:", myRectangle.Area(), "矩形的周长:", myRectangle.Perimeter(),
		"圆的面积:", myCircle.Area(), "圆的周长:", myCircle.Perimeter())

	myEmployee := Employee{Person{"张三", 30}, 1001}
	myEmployee.PrintInfo()
	channelTest1()
	channelTest2()
	MutexTest1()
	MutexTest2()

}

/*
	✅指针

1. 题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
- 考察点 ：指针的使用、值传递与引用传递的区别。
2. 题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
- 考察点 ：指针运算、切片操作。
*/
func pointerTest1(myInt *int) {
	*myInt += 10
	fmt.Println(myInt)
}
func pointerTest2(slicePtr *[]int) {
	// 获取切片引用
	slice := *slicePtr

	fmt.Println("asfasfafffFFf")
	fmt.Println(slicePtr)
	fmt.Println(slice)

	// 遍历切片并修改每个元素
	for i := 0; i < len(slice); i++ {
		slice[i] *= 2
	}
}

/*
## ✅Goroutine
1. 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
  - 考察点 ： go 关键字的使用、协程的并发执行。

2. 题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
  - 考察点 ：协程原理、并发任务调度。
*/
func GoroutineTest1() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Println("偶数:", i)
			}

		}
		wg.Done()
	}()

	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 != 0 {
				fmt.Println("奇数:", i)
			}

		}
		wg.Done()
	}()
	wg.Wait()
}

func GoroutineTest2() {
	// 创建调度器（限制最大并发数为3）
	scheduler := NewTaskScheduler(3)

	// 添加示例任务
	scheduler.AddTask("数据预处理", func() error {
		time.Sleep(800 * time.Millisecond)
		return nil
	})

	scheduler.AddTask("模型训练", func() error {
		time.Sleep(1 * time.Second)
		return nil
	})

	scheduler.AddTask("结果验证", func() error {
		time.Sleep(600 * time.Millisecond)
		return nil
	})

	scheduler.AddTask("报告生成", func() error {
		time.Sleep(400 * time.Millisecond)
		return fmt.Errorf("磁盘空间不足")
	})

	scheduler.AddTask("数据备份", func() error {
		time.Sleep(1 * time.Second)
		return nil
	})

	// 执行任务
	scheduler.Run()

	// 打印详细报告
	scheduler.PrintReport()
}

// Task 定义任务类型
type Task struct {
	ID       int
	Name     string
	Function func() error
}

// TaskResult 存储任务执行结果
type TaskResult struct {
	TaskID    int
	TaskName  string
	Duration  time.Duration
	StartTime time.Time
	EndTime   time.Time
	Err       error
	WorkerID  int
}

// TaskScheduler 任务调度器
type TaskScheduler struct {
	tasks          []Task
	results        []TaskResult
	wg             sync.WaitGroup
	startTime      time.Time
	workerPool     chan struct{}
	mu             sync.Mutex
	completedTasks int
	totalTasks     int
}

// NewTaskScheduler 创建新调度器
func NewTaskScheduler(maxConcurrent int) *TaskScheduler {
	return &TaskScheduler{
		workerPool: make(chan struct{}, maxConcurrent),
		results:    make([]TaskResult, 0),
	}
}

// AddTask 添加任务
func (ts *TaskScheduler) AddTask(name string, taskFunc func() error) {
	ts.tasks = append(ts.tasks, Task{
		ID:       len(ts.tasks),
		Name:     name,
		Function: taskFunc,
	})
}

// Run 执行所有任务
func (ts *TaskScheduler) Run() {
	ts.startTime = time.Now()
	ts.totalTasks = len(ts.tasks)
	ts.completedTasks = 0

	for _, task := range ts.tasks {
		ts.wg.Add(1)
		ts.workerPool <- struct{}{} // 获取worker槽位

		go func(t Task) {
			defer func() {
				<-ts.workerPool // 释放worker槽位
				ts.wg.Done()
			}()

			workerID := len(ts.workerPool)
			ts.executeTask(t, workerID)
		}(task)
	}

	ts.wg.Wait()
}

// 执行单个任务
func (ts *TaskScheduler) executeTask(task Task, workerID int) {
	start := time.Now()
	err := task.Function()
	end := time.Now()

	result := TaskResult{
		TaskID:    task.ID,
		TaskName:  task.Name,
		Duration:  end.Sub(start),
		StartTime: start,
		EndTime:   end,
		Err:       err,
		WorkerID:  workerID,
	}

	ts.mu.Lock()
	ts.results = append(ts.results, result)
	ts.completedTasks++
	ts.mu.Unlock()
}

// PrintReport 打印详细执行报告
func (ts *TaskScheduler) PrintReport() {
	fmt.Println("\n\n=== 任务执行报告 ===")
	fmt.Printf("总任务数: %d\n", ts.totalTasks)
	fmt.Printf("总执行时间: %v\n", time.Since(ts.startTime).Round(time.Millisecond))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\t名称\t状态\t耗时\tWorker\t开始时间\t结束时间")

	for _, res := range ts.results {
		status := "成功"
		if res.Err != nil {
			status = fmt.Sprintf("失败 (%v)", res.Err)
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%v\t%d\t%s\t%s\n",
			res.TaskID,
			res.TaskName,
			status,
			res.Duration.Round(time.Millisecond),
			res.WorkerID,
			res.StartTime.Format("15:04:05.000"),
			res.EndTime.Format("15:04:05.000"),
		)
	}
	w.Flush()
}

/*
## ✅面向对象
1. 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
  - 考察点 ：接口的定义与实现x、面向对象编程风格。
*/
type Shape interface {
	Area() int
	Perimeter() int
}
type Rectangle struct {
	Width  int
	Height int
}

func (r Rectangle) Area() int {
	return r.Width * r.Height
}
func (r Rectangle) Perimeter() int {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius int
}

func (c Circle) Area() int {
	return int(3.14 * float64(c.Radius) * float64(c.Radius))
}
func (c Circle) Perimeter() int {
	return int(2 * 3.14 * float64(c.Radius))
}

/*
## ✅面向对象
2. 题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
  - 考察点 ：组合的使用、方法接收者。
*/

type Person struct {
	Name string
	Age  int
}
type Employee struct {
	Person
	EmployeeID int
}

func (e Employee) PrintInfo() {
	fmt.Printf("员工姓名：%s，年龄：%d，员工ID：%d\n", e.Name, e.Age, e.EmployeeID)
}

/*
## ✅Channel
1. 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
  - 考察点 ：通道的基本使用、协程间通信。
*/
func channelTest1() {
	ch := make(chan int, 10)
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	for num := range ch {
		fmt.Println("通道中的数据：", num)
	}
}

/*
## ✅Channel
2. 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
  - 考察点 ：通道的缓冲机制。
*/
func channelTest2() {
	ch := make(chan int, 10)
	go func() {
		for i := 1; i <= 100; i++ {
			ch <- i
		}
		close(ch)
	}()
	for num := range ch {
		fmt.Println("缓冲通道中的数据：", num)
	}
}

/*
✅锁机制
题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
*/
func MutexTest1() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(10)
	totalNum := 0
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				totalNum++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Println("1.最终的值为：", totalNum)
}

/*
✅锁机制
题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ：原子操作、并发数据安全。
*/
func MutexTest2() {
	var wg sync.WaitGroup
	wg.Add(10)
	var totalNum int32
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt32(&totalNum, 1)
			}
		}()
	}

	wg.Wait()
	fmt.Println("2.最终的值为：", totalNum)
}
