package main

import (
	"fmt"
	"net/http"
	"time"

	"database/sql"
	_ "database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Fish struct {
	Date       string `json:"date"`
	FishName   string `json:"fish_name"`
	Weight     string `json:"weight"`
	Price      string `json:"price"`
	Fraction   string `json:"fraction"`
	Package    string `json:"package"`
	TotalPrice string `json:"total_price"`
}

type FishList struct {
	Fishes []Fish `json:"fishes"`
}
type Customer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
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

	// 建立 today_customer 資料表
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS today_customer (
		ID INTEGER,
		Name TEXT,
		Setting Bool,
		Date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		TotalArrears Float,
		TodayArrears Float
	)`)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("today_customer 資料表建立成功")
}

func insertSelectCustomer(name string, id int, date string) error {

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	now := time.Now()
	date_now := now.Format("2006-01-02")

	_, err = db.Exec(`INSERT INTO today_customer (Name, ID, Setting, Date,TotalArrears,TodayArrears) VALUES (?, ?, 0,?,?,?)`, name, id, date_now, 102.55, 103.55)
	if err != nil {
		return fmt.Errorf("failed to insert customer: %v", err)
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
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// 登入頁面
	router.POST("/login", func(c *gin.Context) {

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

	// 列印頁面
	router.GET("/print", func(c *gin.Context) {
		c.HTML(http.StatusOK, "print.html", gin.H{})
	})

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

	// 測試
	router.POST("/fish", handlePostFish)

	// 取得客戶資訊
	router.GET("/get_customer_name", handleCustome)

	// 選擇今天客戶
	router.POST("/set_today_customer_name", set_today_customer_name)

	// 取得魚種資訊
	router.GET("/get_product_name", handleProduct)

	// 讀取客戶頁面
	router.GET("/select_customer", func(c *gin.Context) {
		c.HTML(http.StatusOK, "select_customer.html", gin.H{})
	})

	router.Run(":8080")
}

func handlePostFish(c *gin.Context) {
	var fishes []Fish
	if err := c.BindJSON(&fishes); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(fishes)

	c.JSON(200, gin.H{"message": "Success"})
}

func handleCustome(c *gin.Context) {
	customers := []Customer{}

	for i := 1; i <= 120; i++ {
		name := fmt.Sprintf("測試員(%d)", i)
		customers = append(customers, Customer{i, name, ""})
	}
	c.JSON(http.StatusOK, customers)
}

func set_today_customer_name(c *gin.Context) {
	var customers []Customer
	if err := c.BindJSON(&customers); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	for _, customer := range customers {
		insertSelectCustomer(customer.Name, customer.ID, customer.Date)
	}

	fmt.Println(customers)
	c.JSON(200, gin.H{"message": "Success"})
}

func handleProduct(c *gin.Context) {
	customers := []Product{
		{1, "白鯧", "01"},
		{2, "黑鯧", "02"},
		{3, "大頭鰱", "03"},
	}
	c.JSON(http.StatusOK, customers)
}
