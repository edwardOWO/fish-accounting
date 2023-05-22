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
	ID           int     `json:"id"`
	Date         string  `json:"date"`
	FishName     string  `json:"fishName"`
	Weight       float32 `json:"weight"`
	Price        int     `json:"price"`
	Fraction     float32 `json:"fraction"`
	Package      string  `json:"package"`
	TotalPrice   int     `json:"totalPrice"`
	CustomerName string  `json:"customerName"`
	INDEX        int     `json:"index"`
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
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Customer (
			ID INTEGER,
			Name TEXT,
			TotalArrears Int
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
		DataIndex INTEGER
	)`)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 建立初始化使用者

	fmt.Println("today_customer 資料表建立成功")
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
		fmt.Print(PaymentsResult)
		_, err = db.Exec(`UPDATE today_customer SET Setting = ?, Sort = ?,PaymentsResult = ? WHERE date = ? AND id = ?;`, setting, sort, PaymentsResult, date, id)
		if err != nil {
			return fmt.Errorf("failed to insert customer: %v", err)
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

	// 輸入帳目資料
	router.POST("/accountDetail", handlePostFish)

	// 取得個人帳目資料
	router.GET("/accountDetail", getCustomAccount)

	// 取得客戶資訊
	router.GET("/get_customer_name", handleCustome)

	// 取得帳目資訊
	router.GET("/get_all_account_customer", getAllAccountCustomer)

	// 還帳目
	router.GET("/clear", clear)

	// 選擇今天客戶
	router.POST("/set_today_customer_name", set_today_customer_name)

	// 讀取今天客戶
	router.GET("/get_today_customer_name", get_today_customer_name)

	// 讀取今天的帳款
	router.GET("/get_customer_account_date", get_customer_account_date)

	// 取得魚種資訊
	router.GET("/get_product_name", handleProduct)

	// 讀取客戶頁面
	router.GET("/select_customer", func(c *gin.Context) {
		c.HTML(http.StatusOK, "select_customer.html", gin.H{})
	})

	router.Run(":8080")
}

func clear(c *gin.Context) {

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	id := c.Query("id")
	date := c.Query("date")
	payment := c.Query("payment")
	clear := c.Query("clear")

	layout := "2006-01-02"

	t, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	if clear == "" {
		// 部分還賬
		_, err = db.Exec(`UPDATE today_customer SET Payments=? WHERE ID=? AND Date=?`, payment, id, t)

		if err != nil {
			fmt.Print(err.Error())
		}
	} else {
		// 完整還賬
		_, err = db.Exec(`UPDATE today_customer SET Clear=? WHERE ID=?`, true, id)

		if err != nil {
			fmt.Print(err.Error())
		}

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
		err := db.QueryRow("SELECT COUNT(*) FROM accountDetail WHERE DataIndex = ? AND Date = ?", detail.INDEX, t).Scan(&count)
		if err != nil {
			count = 1
		}

		if count > 0 {
			// 执行更新操作
			_, err = db.Exec("UPDATE accountDetail SET ID = ?, CustomerName = ?, FishName = ?, Weight = ?, Price = ?, Fraction = ?, Package = ?, TotalPrice = ?, Print = ? WHERE DataIndex = ? AND Date = ?",
				detail.ID, detail.CustomerName, detail.FishName, detail.Weight, detail.Price, detail.Fraction, detail.Package, detail.TotalPrice, false, detail.INDEX, t)
		} else {
			// 执行插入操作
			_, err = db.Exec("INSERT INTO accountDetail (ID, CustomerName, Date, FishName, Weight, Price, Fraction, Package, TotalPrice, Print, DataIndex) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				detail.ID, detail.CustomerName, t, detail.FishName, detail.Weight, detail.Price, detail.Fraction, detail.Package, detail.TotalPrice, false, detail.INDEX)
		}

	}

	// 更新今天的帳目,並且標記已經處理

	// 加總使用者帳款,如果 Clear 欄位為1(true),表示已經還款故不再加總計算
	rows, err := db.Query(`select TotalPrice from  accountDetail WHERE  ID=? and Date=?`, fishes[0].ID, t)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	TotalArrears := 0
	for rows.Next() {
		data := 0
		err := rows.Scan(&data)

		if err != nil {
			log.Fatal(err)
		}

		TotalArrears += data
	}

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

	c.JSON(200, gin.H{"message": "Success"})
}

// 讀取所有客戶
func handleCustome(c *gin.Context) {
	customers := []Customer{}

	db, err := sql.Open("sqlite3", DB_Name)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Customer")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.TotalArrears)
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

		rows, err := db.Query("SELECT * FROM accountDetail where ID =? ORDER BY Date", id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// 迭代查詢結果，並將結果加入 slice
		for rows.Next() {
			var fish Fish
			print := false
			err := rows.Scan(&fish.ID, &fish.CustomerName, &fish.FishName, &fish.Date, &fish.Price, &fish.Weight, &fish.Fraction, &fish.Package, &fish.TotalPrice, &print, &fish.INDEX)
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

		// 迭代查詢結果，並將結果加入 slice
		for rows.Next() {
			var fish Fish
			print := false
			err := rows.Scan(&fish.ID, &fish.CustomerName, &fish.FishName, &fish.Date, &fish.Price, &fish.Weight, &fish.Fraction, &fish.Package, &fish.TotalPrice, &print, &fish.INDEX)
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
