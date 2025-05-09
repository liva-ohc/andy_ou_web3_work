package main

import (
	"fmt"
	"time"
)

const c1, c2, c3 int = 1, 2, 3
const (
	c4 = 4
	c5 = 5
)

func main() {
	fmt.Println("Hello, World!")
	fmt.Println("常量", c1, c2, c3, c4, c5)
	//基础数据类型默认值
	fmt.Println("基础数据类型默认值")
	var a int
	var b string
	var c bool
	fmt.Println(a, b, c)
	//指针
	fmt.Println("指针")
	var d *int = &a
	var e **int = &d
	//打印*d为值，而不是地址
	fmt.Println("打印*d为值，而不是地址")
	fmt.Println(d, e, *d)
	//指针修改值
	fmt.Println("指针修改值")
	*d = 10
	fmt.Println(a, *d)
	//for循环的使用 注意go没有while循环
	fmt.Println("for循环的使用")
	for i := 0; i < 10; i++ {
		fmt.Print(" ", i)
	}
	//类似while的使用
	for i := 0; i < 10; {
		fmt.Print(" ", i)
		i++
	}
	//死循环
	// for {
	// 	fmt.Println("死循环")
	// }

	//结构体的使用
	fmt.Println("\n结构体的使用")
	user1 := User{
		Name: "小明",
		Age:  18,
	}
	user2 := User{"小明", 18, 1, "这是一个匿名字段"}
	fmt.Println(user1, user2)
	user1.Name = "小红"
	fmt.Println(user1, user2)

	//if的使用
	fmt.Println("if的使用")
	if user1.Age < 0 {
		fmt.Println("user1年龄小于0")
	} else {
		fmt.Println("user1年龄大于0")
	}
	//switch的使用
	fmt.Println("switch的使用")
	switch user2.Age {
	case 18:
		fmt.Println("user2年龄是18")
	default:
		fmt.Println("user2年龄不是18")
	}
	//select 的使用
	selectChan1 := make(chan int)
	selectChan2 := make(chan int)
	go func() {
		for i := 0; i < 2; i++ {
			selectChan1 <- i
		}
		close(selectChan1)
	}()
	for i := 0; i < 10; i++ {
		fmt.Println("select的使用")
		select {
		case <-selectChan1:
			fmt.Println("selectChan1来数了")
		case <-selectChan2:
			fmt.Println("selectChan2来数了")
		case <-time.After(2 * time.Second):
			fmt.Println("时间到咯!!")

		}
	}

	//outter的使用、goto的使用
	fmt.Println("outter的使用、goto的使用")
outter:
	for i := 0; i < 3; i++ {
		for j := 3; j < 5; j++ {
			fmt.Println("i=", i, "j=", j)
			if j == 4 {
				break outter
			}
		}
	}

	for i := 0; i < 2; i++ {
		if i == 1 {
			goto gototest
		}
	}

gototest:
	fmt.Println("goto的使用")

	//接口的使用
	fmt.Println("接口的使用")
	myCreditcard := &Creditcard{balance: 0, limit: 100}
	pucharseItem(myCreditcard, 20)

	//数组、切片、map、range的使用
	fmt.Println("数组、切片、map、range的使用")
	var arr [5]int
	arr[0] = 1
	var slice = make([]int, 5)
	slice[0] = 1
	maptest := make(map[string]string, 5)
	maptest["name"] = "小明"
	for arrIndex, arrItem := range arr {
		fmt.Println("arrIndex:", arrIndex, "arrItem:", arrItem)
	}
	for sliceIndex, sliceItem := range slice {
		fmt.Println("sliceIndex:", sliceIndex, "sliceItem:", sliceItem)
	}
	for mapKey, mapValue := range maptest {
		fmt.Println("mapKey:", mapKey, "mapValue:", mapValue)
	}

	//数组转切片
	fmt.Println("数组转切片--共享数据")
	arr2 := [...]int{1, 2, 3}
	slice2 := arr2[:]
	slice2[0] = 10
	fmt.Printf("数组转切片--共享数据, 变更后的slice:%v, 变更后的arr:%v\n", slice2, arr2)
	fmt.Println("数组转切片--非共享数据")
	arr3 := [...]int{1, 2, 3}
	slice3 := make([]int, len(arr3))
	copy(slice3, arr3[:])
	slice3[0] = 10
	fmt.Printf("数组转切片--非共享数据, 变更后的slice:%v, 变更后的arr:%v\n", slice3, arr3)

	//goroutine和channel、select的使用
	fmt.Println("goroutine和channel的使用")
	channel1 := make(chan int)
	channel2 := make(chan int)
	go func() {
		defer close(channel1)
		defer close(channel2)
		for i := 0; i < 5; i++ {
			// time.Sleep(1 * time.Second)
			channel1 <- i
			channel2 <- i
		}
	}()

	// go func() {
	// 	defer close(channel2)
	// 	for i := 0; i < 5; i++ {
	// 		// time.Sleep(1 * time.Second)
	// 		channel2 <- i
	// 	}
	// }()

	for i := 0; i < 100; i++ {
		select {
		case data := <-channel1:
			fmt.Println("channel1来数了:", data)
		case data := <-channel2:
			fmt.Println("channel2来数了:", data)
		case <-time.After(10 * time.Second):
			fmt.Println("时间到咯!!")
		default:
			fmt.Println("等待中")
		}
	}

}

type User struct {
	Name string
	Age  int
	int  //匿名字段
	string
}

// 接口的使用
type payMethod interface {
	Account

	Pay(payAmount int) bool
}

type Account interface {
	Balance
}

type Balance interface {
	GetBalance() int
}

type Creditcard struct {
	balance int //余额
	limit   int //额度
}

func (c *Creditcard) Pay(payAmount int) bool {
	if c.balance+payAmount <= c.limit {
		c.balance += payAmount
		fmt.Println("信用卡支付成功, 支付金额为:", payAmount)
		return true
	}
	fmt.Println("信用卡支付失败, 支付金额为:", payAmount)
	return false
}

func (c *Creditcard) GetBalance() int {
	return c.balance
}

func pucharseItem(payMethod payMethod, itemAccount int) {
	if payMethod.Pay(itemAccount) {
		fmt.Println("购买成功，余额为：", payMethod.GetBalance())
	} else {
		fmt.Println("购买失败")
	}
}
