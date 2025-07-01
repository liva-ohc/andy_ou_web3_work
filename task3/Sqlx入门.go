package main

import (
	_ "gorm.io/driver/mysql"
)

type Employee struct {
	Id         uint   `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Name       string `gorm:"TYPE:VARCHAR(255)"`
	Salary     uint   `gorm:"TYPE:BIGINT"`
	Department string `gorm:"TYPE:VARCHAR(255)"`
}

type Book struct {
	ID     uint    `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Title  string  `gorm:"TYPE:VARCHAR(255)"`
	Author string  `gorm:"TYPE:VARCHAR(255)"`
	Price  float64 `gorm:"TYPE:DECIMAL"`
}

//func main() {
//	//链接数据库
//	dsn := "root:123456@tcp(127.0.0.1:3306)/goTest?charset=utf8mb4&parseTime=True&loc=Local"
//
//	//newLogger := logger.New(
//	//	log.New(os.Stdout, "\r\n", log.LstdFlags),
//	//	logger.Config{
//	//		LogLevel: logger.Info,
//	//	},
//	//)
//
//	db, err := sqlx.Connect("mysql", dsn)
//	if err != nil {
//		panic("failed to connect database")
//	}
//	var employees []Employee
//	//编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
//	err = db.Select(&employees, "select * from employees where department=?", "技术部")
//	if err != nil {
//		log.Fatalf("查询失败: %v", err)
//	}
//	fmt.Println(employees)
//	//编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
//	var highestSalaryEmployee Employee
//	err = db.Get(&highestSalaryEmployee, "select * from employees order by salary desc limit 1")
//	if err != nil {
//		log.Fatalf("查询失败: %v", err)
//	}
//	fmt.Println(highestSalaryEmployee)
//
//	//编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全
//	var books []Book
//	err = db.Select(&books, "select * from books where price > ?", 50)
//	if err != nil {
//		log.Fatalf("查询失败: %v", err)
//	}
//	fmt.Println(books)
//}
