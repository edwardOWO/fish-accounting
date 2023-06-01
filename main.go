package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"database/sql"
	_ "database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Fish struct {
	ID             int     `json:"id"`
	Date           string  `json:"date"`
	FishName       string  `json:"fishName"`
	Weight         float32 `json:"weight"`
	Price          int     `json:"price"`
	Fraction       float32 `json:"fraction"`
	Package        string  `json:"package"`
	TotalPrice     int     `json:"totalPrice"`
	CustomerName   string  `json:"customerName"`
	INDEX          int     `json:"index"`
	PaymentAmount  int     `json:"paymentamount"`
	PaymentsResult string  `json:"paymentsresult"`
	Clear          bool    `json:"clear"`
}

type FishList struct {
	Fishes []Fish `json:"fishes"`
}
type Customer struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Setting        string    `json:"setting"`
	Date           time.Time `json:"date"`
	TotalArrears   int       `json:"totalArrears"`
	TodayArrears   int       `json:"todayArrears"`
	Payments       int       `json:"Payments"`
	PaymentsResult string    `json:"PaymentsResult"`
	Clear          bool      `json:"Clear"`
	Sort           int       `json:"sort"`
	Print          bool      `json:"print"`
}
type CustomerList struct {
	Customers []Customer `json:"customers"`
}

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

var DB_Name = "test.sqlite"

func init_db() {
	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	//testPresure()

	// 建立 today_customer 資料表
	// 未印帳款加總
	// 當前帳目加總
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Customer (
			ID INTEGER,
			Name TEXT,
			TotalArrears Int,
			TodayArrears Int,
			Print Bool
		)`)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 建立 today_customer 資料表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS today_customer (
		ID INTEGER,
		Name TEXT,
		Setting Bool,
		Date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		TotalArrears Float,
		TodayArrears Int,
		Payments Int,
		Clear Bool,
		Sort INTEGER,
		PaymentsResult TEXT
	)`)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 建立 today_customer 資料表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS accountDetail (
		ID INTEGER,
		CustomerName TEXT,
		FishName TEXT,
		Date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		Price INTEGER,
		Weight Float,
		Fraction TEXT,
		Package TEXT,
		TotalPrice INTEGER,
		Print Bool,
		DataIndex INTEGER,
		PaymentsResult Text,
		PaymentAmount INTEGER,
		Clear Bool
	)`)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 建立初始化使用者

	fmt.Println("today_customer 資料表建立成功")
}
func testPresure() {

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	id := 10
	detail := Fish{}
	detail.ID = id
	detail.CustomerName = ""
	detail.Date = "2023-05-01 00:00:00+00:00"
	detail.FishName = ""
	detail.Fraction = float32(0)
	detail.INDEX = 999
	detail.Package = ""
	detail.PaymentAmount = 0
	detail.PaymentsResult = "共:"
	detail.Price = 0
	detail.TotalPrice = 0
	detail.Clear = true

	for i := 0; i < 100000; i++ {
		db.Exec("INSERT INTO accountDetail (ID, CustomerName, Date, FishName, Weight, Price, Fraction, Package, TotalPrice, Print, DataIndex,PaymentsResult,Clear,PaymentAmount ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?)",
			detail.ID, detail.CustomerName, detail.Date, detail.FishName, detail.Weight, detail.Price, detail.Fraction, detail.Package, detail.TotalPrice, true, detail.INDEX, detail.PaymentsResult, detail.Clear, detail.PaymentAmount)

	}

}

func insertSelectCustomer(setting string, name string, id int, date time.Time, sort int, TodayArrears int, PaymentsResult string) error {

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	// 先檢查今天的選擇資料是否存在

	rows, err := db.Query(`select * from  today_customer where name=? and id=? and date=?`, name, id, date)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	index := 0
	// 迭代查詢結果，並將結果加入 slice
	for rows.Next() {
		index++
	}

	if index == 0 {
		_, err = db.Exec(`INSERT INTO today_customer (Name, ID, Setting, Date,TotalArrears,TodayArrears,Sort,Payments,Clear,PaymentsResult) VALUES (?, ?, 0,?,?,?,?,0,false,'')`, name, id, date, 0, 0, sort)
		if err != nil {
			return fmt.Errorf("failed to insert customer: %v", err)
		}
	} else {

		if PaymentsResult != "" {
			_, err = db.Exec(`UPDATE today_customer SET Setting = ?, Sort = ?,PaymentsResult = ? WHERE date = ? AND id = ?;`, setting, sort, PaymentsResult, date, id)
			if err != nil {
				return fmt.Errorf("failed to insert customer: %v", err)
			}
		} else {
			// 當結帳訊息沒有資訊時,不更新結帳訊息
			_, err = db.Exec(`UPDATE today_customer SET Setting = ?, Sort = ? WHERE date = ? AND id = ?;`, setting, sort, date, id)
			if err != nil {
				return fmt.Errorf("failed to insert customer: %v", err)
			}
		}

	}

	return nil
}

func main() {

	// 初始化 DB
	init_db()

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/script", "./script")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "menu.html", gin.H{})
	})

	// 登入頁面
	router.GET("/login", func(c *gin.Context) {

		c.HTML(http.StatusOK, "menu.html", gin.H{})

	})

	// 帳目輸入頁面
	router.GET("/input", func(c *gin.Context) {
		c.HTML(http.StatusOK, "input.html", gin.H{})
	})

	// 帳目檢查頁面
	router.GET("/check", func(c *gin.Context) {
		c.HTML(http.StatusOK, "check.html", gin.H{})
	})

	// 設定列印人頁面
	router.GET("/select_print", func(c *gin.Context) {
		c.HTML(http.StatusOK, "select_print.html", gin.H{})
	})

	// 選擇列印頁面
	router.GET("/print", generatePrintHTML)

	// 產生列印檔案
	router.POST("/print", generatePrintDetail)

	// 帳目狀況
	router.GET("/status", func(c *gin.Context) {
		c.HTML(http.StatusOK, "status.html", gin.H{})
	})

	// 客戶建檔
	router.GET("/customer", func(c *gin.Context) {
		c.HTML(http.StatusOK, "customer.html", gin.H{})
	})

	// 魚名建檔
	router.GET("/product", func(c *gin.Context) {
		c.HTML(http.StatusOK, "product.html", gin.H{})
	})

	// 當日總帳
	router.GET("/account", func(c *gin.Context) {
		c.HTML(http.StatusOK, "account.html", gin.H{})
	})

	// 輸入帳目資料
	router.POST("/accountDetail", handlePostFish)

	// 輸入帳目資料
	router.POST("/payment", payment)

	// 取得個人帳目資料
	router.GET("/accountDetail", getCustomAccount)

	// 取得客戶資訊
	router.GET("/get_customer_name", handleCustome)

	// 取得帳目資訊
	router.GET("/get_all_account_customer", getAllAccountCustomer)

	// 還帳目
	router.POST("/clear", clear)

	// 選擇今天客戶
	router.POST("/set_today_customer_name", set_today_customer_name)

	// 讀取今天客戶
	router.GET("/get_today_customer_name", get_today_customer_name)

	// 選則下一個客戶
	router.POST("/next_customer", next_customer)

	// 選則下一個客戶
	router.POST("/PrintAndClose", PrintAndClose)

	// 更新當前最新帳款
	router.POST("/UpdateTodayArrears", UpdateTodayArrears)

	// 讀取今天的帳款
	router.GET("/get_customer_account_date", get_customer_account_date)

	// 讀取帳目結果
	router.GET("/get_customer_account_result", get_customer_account_result)

	// 取得魚種資訊
	router.GET("/get_product_name", handleProduct)

	// 讀取客戶頁面
	router.GET("/select_customer", func(c *gin.Context) {
		c.HTML(http.StatusOK, "select_customer.html", gin.H{})
	})

	router.Run(":8080")
}

func UpdateTodayArrears(c *gin.Context) {

	id := c.Query("id")

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var fishes []Fish
	if err := c.BindJSON(&fishes); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	rows2, err := db.Query(`select TodayArrears  from  Customer  WHERE  ID=?`, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()

	TodayArrears := 0
	for rows2.Next() {
		rows2.Scan(&TodayArrears)
	}

	// 計算當前欠款金額
	rows, err := db.Query(`select TotalPrice,PaymentAmount from  accountDetail WHERE  ID=? and Clear=false`, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	TotalArrears := 0
	Imcome := 0
	for rows.Next() {
		TotalPrice := 0
		PaymentAmount := 0
		err := rows.Scan(&TotalPrice, &PaymentAmount)

		if err != nil {
			log.Fatal(err)
		}

		TotalArrears += TotalPrice
		Imcome += PaymentAmount
	}
	TotalArrears = TotalArrears - Imcome

	if TodayArrears != TotalArrears && TotalArrears > 0 {

		// 更新當前最新帳款
		_, err = db.Exec("UPDATE Customer SET TodayArrears = ?,Print = ? WHERE ID = ?", TotalArrears, false, id)
		if err != nil {

			fmt.Print(err.Error())
		}

	} else {
		// 更新當前最新帳款
		_, err = db.Exec("UPDATE Customer SET TodayArrears = ? WHERE ID = ?", TotalArrears, id)
		if err != nil {

			fmt.Print(err.Error())
		}
	}

}

func clear(c *gin.Context) {

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var fishes []Fish
	if err := c.BindJSON(&fishes); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"

	t, err := time.Parse(layout, fishes[0].Date)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	//TotalArrears -= fishes[0].PaymentAmount

	for _, detail := range fishes {

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM accountDetail WHERE DataIndex = ? AND Date = ? AND ID = ?", detail.INDEX, t, detail.ID).Scan(&count)
		if err != nil {
			count = 1
		}

		if count > 0 {

			_, err = db.Exec("UPDATE accountDetail SET ID = ?, CustomerName = ?, FishName = ?, Weight = ?, Price = ?, Fraction = ?, Package = ?, TotalPrice = ?, Print = ?,Clear = ?, PaymentsResult = ?,PaymentAmount= ? WHERE DataIndex = ? AND Date = ? AND ID = ?",
				detail.ID, detail.CustomerName, detail.FishName, detail.Weight, detail.Price, detail.Fraction, detail.Package, detail.TotalPrice, false, detail.Clear, "", fishes[0].PaymentAmount, detail.INDEX, t, detail.ID)
			if err != nil {

				fmt.Print(err.Error())
			}
		} else {

			_, err = db.Exec("INSERT INTO accountDetail (ID, CustomerName, Date, FishName, Weight, Price, Fraction, Package, TotalPrice, Print, DataIndex,PaymentsResult,Clear,PaymentAmount ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?)",
				detail.ID, detail.CustomerName, t, detail.FishName, detail.Weight, detail.Price, detail.Fraction, detail.Package, detail.TotalPrice, false, detail.INDEX, "", detail.Clear, fishes[0].PaymentAmount)
			if err != nil {
				fmt.Print(err.Error())
			}
		}

	}

	//rows, err := db.Query(`select TotalPrice from  accountDetail WHERE  ID=? and Date=? and Clear=false`, fishes[0].ID, t)
	rows, err := db.Query(`select TotalPrice,PaymentAmount from  accountDetail WHERE  ID=? and Clear=false`, fishes[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	TotalArrears := 0
	Imcome := 0
	for rows.Next() {
		TotalPrice := 0
		PaymentAmount := 0
		err := rows.Scan(&TotalPrice, &PaymentAmount)

		if err != nil {
			log.Fatal(err)
		}

		TotalArrears += TotalPrice
		Imcome += PaymentAmount
	}

	result := ""
	if fishes[0].PaymentAmount != 0 {
		result += t.Format("01/02")

		if TotalArrears-Imcome >= 0 {
			result += " 入=" + strconv.Itoa(fishes[0].PaymentAmount) + " 欠=" + strconv.Itoa(TotalArrears-Imcome)
		} else {
			result += " 入=" + strconv.Itoa(fishes[0].PaymentAmount) + " 剩=" + strconv.Itoa(TotalArrears-Imcome)
		}

	}

	_, err = db.Exec("UPDATE accountDetail SET PaymentsResult= ? WHERE DataIndex = ? AND Date = ?", result, fishes[0].INDEX, t)
	if err != nil {
		fmt.Print(err.Error())
	}

}
func PrintAndClose(c *gin.Context) {

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var fishes []Fish
	if err := c.BindJSON(&fishes); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"

	t, err := time.Parse(layout, fishes[0].Date)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	//TotalArrears -= fishes[0].PaymentAmount

	for _, detail := range fishes {

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM accountDetail WHERE DataIndex = ? AND Date = ?", detail.INDEX, t).Scan(&count)
		if err != nil {
			count = 1
		}

		if count > 0 {

			_, err = db.Exec("UPDATE accountDetail SET ID = ?, CustomerName = ?, FishName = ?, Weight = ?, Price = ?, Fraction = ?, Package = ?, TotalPrice = ?, Print = ?,Clear = ?, PaymentsResult = ?,PaymentAmount= ? WHERE DataIndex = ? AND Date = ?",
				detail.ID, detail.CustomerName, detail.FishName, detail.Weight, detail.Price, detail.Fraction, detail.Package, detail.TotalPrice, false, detail.Clear, "", fishes[0].PaymentAmount, detail.INDEX, t)
			if err != nil {

				fmt.Print(err.Error())
			}
		} else {

			_, err = db.Exec("INSERT INTO accountDetail (ID, CustomerName, Date, FishName, Weight, Price, Fraction, Package, TotalPrice, Print, DataIndex,PaymentsResult,Clear,PaymentAmount ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?)",
				detail.ID, detail.CustomerName, t, detail.FishName, detail.Weight, detail.Price, detail.Fraction, detail.Package, detail.TotalPrice, false, detail.INDEX, "", detail.Clear, fishes[0].PaymentAmount)
			if err != nil {
				fmt.Print(err.Error())
			}
		}

	}

	//rows, err := db.Query(`select TotalPrice from  accountDetail WHERE  ID=? and Date=? and Clear=false`, fishes[0].ID, t)
	rows, err := db.Query(`select TotalPrice,PaymentAmount from  accountDetail WHERE  ID=? and Clear=false`, fishes[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	TotalArrears := 0
	Imcome := 0
	for rows.Next() {
		TotalPrice := 0
		PaymentAmount := 0
		err := rows.Scan(&TotalPrice, &PaymentAmount)

		if err != nil {
			log.Fatal(err)
		}

		TotalArrears += TotalPrice
		Imcome += PaymentAmount
	}
	result := "共="
	result += strconv.Itoa(TotalArrears - Imcome)

	// 結帳
	_, err = db.Exec("UPDATE accountDetail SET PaymentsResult= ? WHERE DataIndex = ? AND Date = ?", result, fishes[0].INDEX, t)
	if err != nil {
		fmt.Print(err.Error())
	}

	// 更新前帳金額
	_, err = db.Exec("UPDATE Customer SET TotalArrears = ? WHERE ID = ?", TotalArrears-Imcome, fishes[0].ID)
	if err != nil {

		fmt.Print(err.Error())
	}

	if TotalArrears-Imcome == 0 {

		// 將單子設定成已經結帳
		_, err = db.Exec("UPDATE accountDetail SET Clear = ? WHERE ID = ?", true, fishes[0].ID)
		if err != nil {

			fmt.Print(err.Error())
		}
	}

	// 將單子設定成已印出
	_, err = db.Exec("UPDATE accountDetail SET Print = ? WHERE ID = ?", true, fishes[0].ID)
	if err != nil {

		fmt.Print(err.Error())
	}

}

// 選擇下一個使用者
func next_customer(c *gin.Context) {

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var fishes []Fish
	if err := c.BindJSON(&fishes); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"

	t, err := time.Parse(layout, fishes[0].Date)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}
	_, err = db.Exec("UPDATE today_customer SET Setting = ? WHERE ID = ? AND Date = ?", true, fishes[0].ID, t)
	if err != nil {

		fmt.Print(err.Error())
	}
}

func handlePostFish(c *gin.Context) {
	// 更新客戶資料
	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// 更新今日詳細帳目
	var fishes []Fish
	if err := c.BindJSON(&fishes); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	// 解析日期字符串

	day := ""
	day = fishes[0].Date

	t, err := time.Parse(layout, day)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	} // 解析日期字符串

	//  先刪除當天的所有數據,待後續寫入數據
	//_, err = db.Exec(`DELETE from accountDetail WHERE date(Date) = date(?) AND ID=?`, t, fishes[0].ID)

	//if err != nil {
	//	fmt.Print(err.Error())
	//}

	// 寫入當天的所有數據
	for _, detail := range fishes {

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM accountDetail WHERE DataIndex = ? AND Date = ? AND ID=?", detail.INDEX, t, fishes[0].ID).Scan(&count)
		if err != nil {
			count = 1
		}

		if count > 0 {
			// 执行更新操作
			_, err = db.Exec("UPDATE accountDetail SET ID = ?, CustomerName = ?, FishName = ?, Weight = ?, Price = ?, Fraction = ?, Package = ?, TotalPrice = ?, Print = ?,Clear = ?, PaymentsResult = ? ,PaymentAmount = ? WHERE DataIndex = ? AND Date = ? AND ID=?",
				detail.ID, detail.CustomerName, detail.FishName, detail.Weight, detail.Price, detail.Fraction, detail.Package, detail.TotalPrice, false, detail.Clear, detail.PaymentsResult, detail.PaymentAmount, detail.INDEX, t, fishes[0].ID)
		} else {
			// 执行插入操作
			_, err = db.Exec("INSERT INTO accountDetail (ID, CustomerName, Date, FishName, Weight, Price, Fraction, Package, TotalPrice, Print, DataIndex,PaymentsResult,Clear, PaymentAmount) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,0)",
				detail.ID, detail.CustomerName, t, detail.FishName, detail.Weight, detail.Price, detail.Fraction, detail.Package, detail.TotalPrice, false, detail.INDEX, detail.PaymentsResult, detail.Clear, detail.PaymentAmount)
			if err != nil {
				fmt.Print(err.Error())
			}
		}

	}

	c.JSON(200, gin.H{"message": "Success"})
}
func payment(c *gin.Context) {

	/*
		query := "UPDATE today_customer SET Setting=?,TodayArrears=? WHERE date=? AND ID=?"
		result, err := db.Exec(query, 0, TotalArrears, t, fishes[0].ID)
		fmt.Print(result)
		if err != nil {
			log.Fatal(err)
		}

		// 加總使用者帳款,如果 Clear 欄位為1(true),表示已經還款故不再加總計算
		TotalArrears = 0
		rows, err = db.Query(`select TodayArrears from  today_customer WHERE  ID=? and Clear=0`, fishes[0].ID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			data := 0
			err := rows.Scan(&data)

			if err != nil {
				log.Fatal(err)
			}

			TotalArrears += data
		}

		// 更新當前所有的欠款數到 Customer
		query = "UPDATE Customer SET TotalArrears=? where id=?"
		_, err = db.Exec(query, TotalArrears, fishes[0].ID)

		if err != nil {
			log.Fatal(err)
		}
	*/
}

// 讀取所有客戶
func handleCustome(c *gin.Context) {
	customers := []Customer{}

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT ID,Name,TotalArrears,TodayArrears,Print FROM Customer")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.TotalArrears, &customer.TodayArrears, &customer.Print)
		if err != nil {
			log.Fatal(err)
		}
		customers = append(customers, customer)
	}
	c.JSON(http.StatusOK, customers)
}

func getAllAccountCustomer(c *gin.Context) {

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var getAllAccountCustomer []Customer

	rows, err := db.Query("SELECT * FROM today_customer where Setting =1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 迭代查詢結果，並將結果加入 slice
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Setting, &customer.Date, &customer.TodayArrears, &customer.TotalArrears, &customer.Sort, &customer.Payments, &customer.PaymentsResult)
		if err != nil {
			log.Fatal(err)
		}
		getAllAccountCustomer = append(getAllAccountCustomer, customer)
	}

	// 檢查是否有迭代中發生錯誤
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, getAllAccountCustomer)

}

func getCustomAccount(c *gin.Context) {

	id := c.Query("id")
	date := c.Query("date")

	fmt.Print(date)
	// 在这里使用 id 参数进行逻辑处理
	// ...

	if date == "" && id != "" {

		db, err := sql.Open("sqlite3", DB_Name)
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()

		var getCustomAccountDetail []Fish

		rows, err := db.Query("SELECT * FROM accountDetail where ID =? ORDER BY Date,DataIndex;", id)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// 迭代查詢結果，並將結果加入 slice
		for rows.Next() {
			var fish Fish
			print := false
			err := rows.Scan(&fish.ID, &fish.CustomerName, &fish.FishName, &fish.Date, &fish.Price, &fish.Weight, &fish.Fraction, &fish.Package, &fish.TotalPrice, &print, &fish.INDEX, &fish.PaymentsResult, &fish.PaymentAmount, &fish.Clear)
			if err != nil {
				log.Fatal(err)
			}
			getCustomAccountDetail = append(getCustomAccountDetail, fish)
		}

		// 檢查是否有迭代中發生錯誤
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, getCustomAccountDetail)

	} else {
		db, err := sql.Open("sqlite3", DB_Name)
		if err != nil {
			fmt.Println(err)
		}
		defer db.Close()

		var getCustomAccountDetail []Fish

		dateString := date
		dateLayout := "2006-01-02" // 指定日期字符串的格式

		date, err := time.Parse(dateLayout, dateString)
		if err != nil {
			fmt.Println("日期解析错误:", err)
			return
		}

		rows, err := db.Query("SELECT * FROM accountDetail WHERE date(Date) = date(?) AND ID=?", date, id)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var fish Fish
			print := false
			err := rows.Scan(&fish.ID, &fish.CustomerName, &fish.FishName, &fish.Date, &fish.Price, &fish.Weight, &fish.Fraction, &fish.Package, &fish.TotalPrice, &print, &fish.INDEX, &fish.PaymentAmount)
			if err != nil {
				log.Fatal(err)
			}
			getCustomAccountDetail = append(getCustomAccountDetail, fish)
		}

		// 檢查是否有迭代中發生錯誤
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, getCustomAccountDetail)

	}

}

func set_today_customer_name(c *gin.Context) {
	var customers []Customer
	if err := c.BindJSON(&customers); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for _, customer := range customers {
		insertSelectCustomer(customer.Setting, customer.Name, customer.ID, customer.Date, customer.Sort, customer.TodayArrears, customer.PaymentsResult)
	}

	fmt.Println(customers)
	c.JSON(200, gin.H{"message": "Success"})
}

func get_customer_account_date(c *gin.Context) {
	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	id := c.Query("id")

	rows, err := db.Query("SELECT Date,Clear FROM today_customer where ID =? ORDER BY Date DESC", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 迭代查詢結果，並將結果加入 slice
	var get_customer_account_date []string
	for rows.Next() {

		//t := time.Date(2023, time.May, 20, 0, 0, 0, 0, time.UTC)

		//str := t.Format("2006-01-02")

		str := ""

		// 檢查還款狀態
		clear_status := false
		err := rows.Scan(&str, &clear_status)
		// 解析时间字符串为 time.Time 类型
		t, err := time.Parse(time.RFC3339, str)
		if err != nil {
			fmt.Println("时间解析错误:", err)
			return
		}

		// 将 time.Time 格式化为 "YYYYMMDD" 格式的字符串
		formattedStr := t.Format("2006-01-02")
		if err != nil {
			log.Fatal(err)
		}
		formattedStr += ","
		clear_status_ := strconv.FormatBool(clear_status)
		formattedStr += clear_status_
		get_customer_account_date = append(get_customer_account_date, formattedStr)
	}

	// 檢查是否有迭代中發生錯誤
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, get_customer_account_date)

}

func get_customer_account_result(c *gin.Context) {
	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	id := c.Query("id")
	date := c.Query("date")

	dateString := date
	dateLayout := "2006-01-02" // 指定日期字符串的格式

	t, err := time.Parse(dateLayout, dateString)
	if err != nil {
		fmt.Println("日期解析错误:", err)
		return
	}

	rows, err := db.Query("SELECT PaymentsResult  FROM today_customer where ID =? AND Date=? ORDER BY Date DESC", id, t)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var get_customer_account_result []string
	for rows.Next() {

		str := ""
		err := rows.Scan(&str)

		if err != nil {
			log.Fatal(err)
		}

		get_customer_account_result = append(get_customer_account_result, str)
	}

	// 檢查是否有迭代中發生錯誤
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, get_customer_account_result)

}

func get_today_customer_name(c *gin.Context) {

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var todayCustomers []Customer

	rows, err := db.Query("SELECT * FROM today_customer WHERE Setting=0 ORDER BY Sort ASC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 迭代查詢結果，並將結果加入 slice
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Setting, &customer.Date, &customer.TotalArrears, &customer.TodayArrears, &customer.Payments, &customer.Clear, &customer.Sort, &customer.PaymentsResult)
		if err != nil {
			log.Fatal(err)
		}
		todayCustomers = append(todayCustomers, customer)
	}

	// 檢查是否有迭代中發生錯誤
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, todayCustomers)

	// 輸出結果
	fmt.Println(todayCustomers)
}

func handleProduct(c *gin.Context) {
	Products := []Product{
		{1, "白鯧", "0"},
		{2, "黑鯧", "1"},
		{3, "銀鯧", "2"},
		{4, "刺鯧", "3"},
		{5, "尊魚", "4"},
		{6, "四破", "5"},
		{7, "花輝", "6"},
		{8, "紅尾", "7"},
		{9, "黑尾", "8"},
		{10, "盤仔", "9"},
		{11, "赤宗", "10"},
		{12, "赤羽", "11"},
		{13, "赤筆", "12"},
		{14, "赤海", "13"},
		{15, "赤目", "14"},
		{16, "黑格", "15"},
		{17, "刀", "16"},
		{18, "金線", "17"},
		{19, "火口", "18"},
		{20, "黃花", "19"},
		{21, "白北", "20"},
		{22, "馬加", "21"},
		{23, "七爐", "22"},
		{24, "午", "23"},
		{25, "白口", "24"},
		{26, "花枝", "25"},
		{27, "支肉", "26"},
		{28, "軟支", "27"},
		{29, "章", "28"},
		{30, "卷", "29"},
		{31, "熟卷", "30"},
		{32, "尤", "31"},
		{33, "蝦", "32"},
		{34, "蝦仁", "33"},
		{35, "市", "34"},
		{36, "市足", "35"},
		{37, "瓜子", "36"},
		{38, "蛤", "37"},
		{39, "竹蛤", "38"},
		{40, "螺", "39"},
		{41, "草魚", "40"},
		{42, "連魚", "41"},
		{43, "南代", "42"},
		{44, "虱目", "43"},
		{45, "虱頭", "44"},
		{46, "虱肚", "45"},
		{47, "烏", "46"},
		{48, "秋刀", "47"},
		{49, "加納", "48"},
		{50, "龍尖", "49"},
		{51, "皮刀", "50"},
		{52, "土托", "51"},
		{53, "鐵甲", "52"},
		{54, "飛魚", "53"},
		{55, "肉魚", "54"},
		{56, "兔魚", "55"},
		{57, "沙腸", "56"},
		{58, "石喬", "57"},
		{59, "油魚", "58"},
		{60, "雪魚", "59"},
		{61, "紅冬", "60"},
		{62, "甘魚", "61"},
		{63, "三紋", "62"},
		{64, "紅條", "63"},
		{65, "三牙", "64"},
		{66, "花身", "65"},
		{67, "平瓜", "66"},
		{68, "生", "67"},
		{69, "串", "68"},
		{70, "蟳", "69"},
		{71, "扁魚", "70"},
		{72, "溫", "71"},
		{73, "石斑", "72"},
		{74, "秋姑", "73"},
		{75, "力魚", "74"},
		{76, "紅槽", "75"},
		{77, "柳葉", "76"},
		{78, "象耳", "77"},
		{79, "鹹卷", "78"},
		{80, "蚵", "79"},
		{81, "長加", "80"},
		{82, "狗母", "81"},
		{83, "方", "82"},
		{84, "紅魚", "83"},
		{85, "白松", "84"},
		{86, "銀花", "85"},
		{87, "補", "86"},
		{88, "丁香", "87"},
		{89, "紅衫", "88"},
		{90, "尖梭", "89"},
		{91, "鰻魚", "90"},
		{92, "馬頭", "91"},
		{93, "金龍", "92"},
		{94, "雜魚", "93"},
		{95, "皮魚", "94"},
		{96, "奎魚", "95"},
		{97, "破北", "96"},
		{98, "煙虎", "97"},
		{99, "厚唇", "98"},
		{100, "牛舌", "99"},
	}
	c.JSON(http.StatusOK, Products)
}
func WriteToFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 寫入資料
	_, err = file.WriteString(data + "\n")
	if err != nil {
		return err
	}

	return nil
}

func fix_print() {

	file, err := os.Open("data2.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	index := 1
	for scanner.Scan() {

		/*
			if index%60 == 0 {
				WriteToFile("data.txt", "")
				WriteToFile("data.txt", "")
				WriteToFile("data.txt", "")
				index = index + 1
			}
		*/
		line := scanner.Text()
		WriteToFile("data.txt", strconv.Itoa(index)+"."+line)
		index = index + 1

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	/*
		for i := 1; i < 120; i++ {

			if i%60 == 0 {
				WriteToFile("data.txt", "")
				WriteToFile("data.txt", "")
				WriteToFile("data.txt", "")
				continue
			}

			if i == 1 {
				WriteToFile("data.txt", "###############################")
			} else if i == 2 {
				WriteToFile("data.txt", "###############################")
			} else {
				WriteToFile("data.txt", "<<"+strconv.Itoa(i))
			}
		}
	*/
}

func generatePrintHTML(c *gin.Context) {

	os.Truncate("templates/print.html", 0)
	// 讀取 txt 檔案內容
	content, err := ioutil.ReadFile("fish.txt")
	if err != nil {
		log.Fatal(err)
	}
	WriteToFile("templates/print.html", string(content))
	// 定義模板
	tmpl := template.Must(template.ParseFiles("templates/print.html"))

	// 將資料傳遞給模板
	data := struct {
		Content string
	}{
		Content: string(content),
	}

	// 產生 HTML
	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		log.Fatal(err)
	}
}

func generatePrintDetail(c *gin.Context) {

	os.Truncate("fish.txt", 0)

	html := `
	<!DOCTYPE html>
<html>
<head>
  <title>列印純文字內容</title>
  <style>
    .text-content {
      white-space: pre;
      font-size: 15px;
	  margin-top: -15px; /* 调整顶部边距 */
    }
	.shorten-distance {
		margin-top: -20px; /* 调整顶部边距 */
		margin-bottom: 5px; /* 调整底部边距 */
		font-size: 12px; /* 缩小字体大小 */
	}
	.shorten-distance2 {
		margin-top: -6px; /* 调整顶部边距 */
	}
  </style>
  </style>
</head>
<body>
	`
	WriteToFile("fish.txt", html)

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	var customers []Customer
	if err := c.BindJSON(&customers); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id := 0
	sum_index := 0
	for _, customer := range customers {

		rows2, err2 := db.Query("SELECT TotalArrears  from Customer where ID=?", customer.ID)
		id = customer.ID
		if err2 != nil {
			log.Fatal(err)
		}
		defer rows2.Close()

		preCount := 0
		for rows2.Next() {
			err = rows2.Scan(&preCount)
		}

		// rows, err := db.Query("SELECT CustomerName,FishName,Date,Price,Weight,Fraction,Package,TotalPrice,PaymentsResult,PaymentAmount,Print  FROM accountDetail WHERE ID=? AND clear=false", customer.ID)

		rows, err := db.Query("SELECT CustomerName,FishName,Date,Price,Weight,Fraction,Package,TotalPrice,PaymentsResult,PaymentAmount,Print  FROM accountDetail WHERE ID=? AND clear=? ORDER BY Date, DataIndex ASC", customer.ID, false)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var fish Fish

		index := 0
		TotalArrears := 0
		Income := 0

		for rows.Next() {
			print := false
			err = rows.Scan(&fish.CustomerName, &fish.FishName, &fish.Date, &fish.Price, &fish.Weight, &fish.Fraction, &fish.Package, &fish.TotalPrice, &fish.PaymentsResult, &fish.PaymentAmount, &print)

			TotalArrears += fish.TotalPrice
			Income += fish.PaymentAmount

			if index == 0 {
				WriteToFile("fish.txt", "<pre class=\"shorten-distance\">-----------------------------------------------</pre>")
				WriteToFile("fish.txt", "<h2 class=\"shorten-distance2\">")
				WriteToFile("fish.txt", fish.CustomerName)
				WriteToFile("fish.txt", "</h2>")
				WriteToFile("fish.txt", "<pre class=\"text-content\">")

				if preCount != 0 {
					WriteToFile("fish.txt", "前帳: "+strconv.Itoa(preCount))
					index++
				}

			}

			// 打印未印的表單
			if print == false {

				if fish.PaymentsResult != "" {
					index++
					WriteToFile("fish.txt", fish.PaymentsResult)
				} else {
					date, err := time.Parse(time.RFC3339, fish.Date)
					if err != nil {
						fmt.Println("日期解析失败:", err)
						return
					}

					// 格式化日期为 MM/DD
					format := "01/02"
					formattedDate := date.Format(format)

					paddedStr := fmt.Sprintf("%-5s %-3s %-6s %-4s %-3s %-4s", formattedDate, fish.FishName, strconv.FormatFloat(float64(fish.Weight), 'f', -1, 32)+"k", strconv.Itoa(fish.Price), strconv.FormatFloat(float64(fish.Fraction), 'f', -1, 32), strconv.Itoa(fish.TotalPrice))
					index++
					WriteToFile("fish.txt", paddedStr)
				}
			}

		}

		if index <= 14 {

			for i := 1; i <= 14-index; i++ {
				WriteToFile("fish.txt", "")
			}
		}

		WriteToFile("fish.txt", "</pre>")
		WriteToFile("fish.txt", "<h4>盛: "+strconv.Itoa(TotalArrears-Income)+"</h4>")
		WriteToFile("fish.txt", "<pre class=\"shorten-distance\">-----------------------------------------------</pre>")

		_, err = db.Exec("UPDATE Customer SET TotalArrears = ? WHERE ID = ?", strconv.Itoa(TotalArrears-Income), id)
		if err != nil {

			fmt.Print(err.Error())
		}

		_, err = db.Exec("UPDATE Customer SET Print = ? WHERE ID = ?", true, id)
		if err != nil {

			fmt.Print(err.Error())
		}

		//_, err = db.Exec("UPDATE Customer SET TodayArrears = ? WHERE ID = ?", 0, id)
		//if err != nil {

		//			fmt.Print(err.Error())
		//}

		if TotalArrears-Income == 0 {

			// 將單子設定成已經結帳
			_, err = db.Exec("UPDATE accountDetail SET Clear = ?,Print = ? WHERE ID = ?", true, true, id)
			if err != nil {

				fmt.Print(err.Error())
			}
		} else {
			_, err = db.Exec("UPDATE accountDetail SET Print = ? WHERE ID = ?", true, id)
			if err != nil {

				fmt.Print(err.Error())
			}
		}

		now := customer.Date

		targetTime := time.Date(2023, time.May, 3, 0, 0, 0, 0, time.UTC)
		targetDateTime := time.Date(now.Year(), now.Month(), now.Day(), targetTime.Hour(), targetTime.Minute(), targetTime.Second(), targetTime.Nanosecond(), time.UTC)
		formattedTime := targetDateTime.Format("2006-01-02 15:04:05-07:00")

		if err != nil {
			fmt.Println("解析错误:", err)
			return
		}

		detail := Fish{}
		detail.ID = id
		detail.CustomerName = ""
		detail.Date = formattedTime
		detail.FishName = ""
		detail.Fraction = float32(0)
		detail.INDEX = 999
		detail.Package = ""
		detail.PaymentAmount = 0
		detail.PaymentsResult = "共:" + strconv.Itoa(TotalArrears-Income)
		detail.Price = 0
		detail.TotalPrice = 0
		detail.Clear = true

		db.Exec("INSERT INTO accountDetail (ID, CustomerName, Date, FishName, Weight, Price, Fraction, Package, TotalPrice, Print, DataIndex,PaymentsResult,Clear,PaymentAmount ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?,?,?)",
			detail.ID, detail.CustomerName, detail.Date, detail.FishName, detail.Weight, detail.Price, detail.Fraction, detail.Package, detail.TotalPrice, true, detail.INDEX, detail.PaymentsResult, detail.Clear, detail.PaymentAmount)
		if err != nil {
			fmt.Print(err.Error())
		}

		sum_index += index

		if sum_index > 29 {
			sum_index = 0
			WriteToFile("fish.txt", "<div style=\"page-break-before: always;\"></div>")
		}

	}

	html_end := `
	<button onclick="window.print()">列印文字</button>
	</body>
	</html>`
	WriteToFile("fish.txt", html_end)

	fmt.Println(customers)
	c.JSON(200, gin.H{"message": "Success"})
}
