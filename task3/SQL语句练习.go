package main

type Student struct {
	Id    uint   `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Name  string `gorm:"TYPE:VARCHAR(255)"`
	Age   int    `gorm:"TYPE:INT"`
	Grade string `gorm:"TYPE:VARCHAR(255)"`
}

type Account struct {
	Id      uint   `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Name    string `gorm:"TYPE:VARCHAR(255)"`
	Balance uint64 `gorm:"TYPE:BIGINT"`
}
type Transaction struct {
	Id            uint `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	FromAccountId uint `gorm:"TYPE:BIGINT"`
	ToAccountId   uint `gorm:"TYPE:BIGINT"`
	Amount        uint `gorm:"TYPE:BIGINT"`
}

//func main() {
//	//链接数据库
//	dsn := "root:123456@tcp(127.0.0.1:3306)/goTest?charset=utf8mb4&parseTime=True&loc=Local"
//
//	newLogger := logger.New(
//		log.New(os.Stdout, "\r\n", log.LstdFlags),
//		logger.Config{
//			LogLevel: logger.Info,
//		},
//	)
//
//	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
//		Logger: newLogger,
//	})
//	if err != nil {
//		panic("failed to connect database")
//	}
//	//创建表
//	migrateErr := db.AutoMigrate(&Student{}, &Account{}, &Transaction{})
//	if migrateErr != nil {
//		log.Println("migrate Err:", migrateErr)
//		return
//	}
//
//	//var students []Student
//	////编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"
//	//db.Create(&Student{Name: "张三", Age: 20, Grade: "三年级"})
//	//
//	////编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息
//	//db.Where("age > ?", 18).Find(&students)
//	//fmt.Println(students)
//	//
//	////编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"
//	//db.Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级")
//	//
//	////编写SQL语句删除 students 表中年龄小于 15 岁的学生记录
//	//db.Where("age < ?", 15).Delete(&Student{})
//
//	//编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，
//	//向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
//
//	//db.Create(&Account{
//	//	Name:    "A",
//	//	Balance: 100,
//	//})
//	//db.Create(&Account{
//	//	Name:    "B",
//	//	Balance: 0,
//	//})
//
//	db.Transaction(func(tx *gorm.DB) error {
//		var accountA Account
//		tx.Where("name = ?", "A").First(&accountA)
//
//		var accountB Account
//		tx.Where("name = ?", "B").First(&accountB)
//		if accountA.Balance < 100 {
//			return errors.New("余额不足")
//		}
//		tx.Where("name = ?", "A").Update("balance", accountA.Balance-100)
//		tx.Where("name = ?", "B").Update("balance", accountA.Balance+100)
//
//		tx.Create(&Transaction{
//			FromAccountId: accountA.Id,
//			ToAccountId:   accountB.Id,
//			Amount:        100,
//		})
//		return nil
//	})
//
//}
