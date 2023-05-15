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
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Setting      string `json:"setting"`
	Date         string `json:"date"`
	TotalArrears string `json:"totalArrears"`
	TodayArrears string `json:"todayArrears"`
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

	_, err = db.Exec(`INSERT INTO today_customer (Name, ID, Setting, Date,TotalArrears,TodayArrears) VALUES (?, ?, 0,?,?,?)`, name, id, date_now, 9000, 10000)
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
	router.GET("/print", generatePrintHTML)

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

	// 讀取今天客戶
	router.GET("/get_today_customer_name", get_today_customer_name)

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

	for i := 1; i <= 30; i++ {
		name := fmt.Sprintf("測試員(%d)", i)
		customers = append(customers, Customer{i, name, "", "", "", ""})
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

func get_today_customer_name(c *gin.Context) {

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	var todayCustomers []Customer

	// 取得今天日期
	today := time.Now().Format("2006-01-02")

	// 查詢今天的 today_customer 資料
	rows, err := db.Query("SELECT * FROM today_customer WHERE date=? ORDER BY ID LIMIT 1", today)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 迭代查詢結果，並將結果加入 slice
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Setting, &customer.Date, &customer.TodayArrears, &customer.TotalArrears)
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
	customers := []Product{
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
	c.JSON(http.StatusOK, customers)
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

	fix_print()
	// 讀取 txt 檔案內容
	content, err := ioutil.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}

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
