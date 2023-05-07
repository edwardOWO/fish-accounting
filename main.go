package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
}
type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

func main() {
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

	// 取得魚種資訊
	router.GET("/get_product_name", handleProduct)

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
	customers := []Customer{
		{1, "王小明"},
		{2, "王曉華"},
		{3, "王小松"},
	}
	c.JSON(http.StatusOK, customers)
}

func handleProduct(c *gin.Context) {
	customers := []Product{
		{1, "白鯧", "01"},
		{2, "黑鯧", "02"},
		{3, "大頭鰱", "03"},
	}
	c.JSON(http.StatusOK, customers)
}
