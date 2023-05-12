let dictionary = {};

window.onload = function() {
    var table = document.getElementById("myTable");
    table.rows[1].cells[0].focus();

    // 宣告 dictionary 變數
    var customer = document.getElementById("customer");
    customer.innerHTML="NULL";

    var current_count = document.getElementById("current_count");
    current_count.innerHTML=0;

    var pre_count = document.getElementById("pre_count");
    pre_count.innerHTML=0;



    // 使用 fetch() 方法呼叫 API
    fetch('/get_product_name')
    .then(response => response.json()) // 將回傳的資料轉為 JSON 格式
    .then(data => {
        // 將 JSON 格式資料存入 dictionary 變數
        for (let i = 0; i < data.length; i++) {
        let customer = data[i];
        dictionary[customer.key] = customer.name;
        }

        // 顯示 dictionary 變數內容
        console.log(dictionary);
    })
    .catch(error => console.error(error)); // 若發生錯誤則顯示錯誤訊息


    // 使用 fetch() 方法呼叫 API
    fetch('/get_today_customer_name')
    .then(response => response.json()) // 將回傳的資料轉為 JSON 格式
    .then(data => {
        // 將 JSON 格式資料存入 dictionary 變數
        for (let i = 0; i < data.length; i++) {
            customer.innerHTML=data[i].name
            pre_count.innerHTML=parseFloat(data[i].totalArrears)
        }

        // 顯示 dictionary 變數內容
        console.log(dictionary);
    })
    .catch(error => console.error(error)); // 若發生錯誤則顯示錯誤訊息


};

myButton.onclick = function() {
    const form = document.getElementById('myTable');
    const tbody = table.getElementsByTagName('tbody')[0];
    const rows = tbody.getElementsByTagName('tr');
    const data = [];
    for (let i = 0; i < rows.length; i++) {
        const row = rows[i];
        const cells = row.getElementsByTagName('td');

        if (cells[0].innerText ==""){
        continue
        }
        const date = cells[0].innerText;
        const fish_name = cells[1].innerText;
        const weight = cells[2].innerText;
        const price = cells[3].innerText;
        const fraction = cells[4].innerText;
        const package = cells[5].innerText;
        const total_price= cells[6].innerText;

        data.push({
        date: date,
        fish_name: fish_name,
        weight: weight,
        price: price,
        fraction: fraction,
        package: package,
        total_price: total_price,
        });
    }
    const xhr = new XMLHttpRequest();
    xhr.open('POST', '/fish');
    xhr.setRequestHeader('Content-Type', 'application/json;charset=UTF-8');
    xhr.send(JSON.stringify(data));
};




function checkFishColor(number) {

    return dictionary[number]
}

var table = document.getElementById("myTable");
var currentRow = 1;
var currentCol = 0;

table.rows[currentRow].cells[currentCol].focus();

table.addEventListener("keydown", function(event) {

    var key = event.which || event.keyCode;

    switch (key) {
        case 37: // left arrow
        if (currentCol > 0) {
            currentCol--;
            table.rows[currentRow].cells[currentCol].focus();
        }
        break;
        case 38: // up arrow
        if (currentRow > 0) {
            currentRow--;
            table.rows[currentRow].cells[currentCol].focus();
        }
        break;
        case 39: // right arrow
        if (currentCol < table.rows[currentRow].cells.length - 1) {
            currentCol++;
            table.rows[currentRow].cells[currentCol].focus();
        }
        break;
        case 40: // down arrow
        if (currentRow < table.rows.length - 1) {
            currentRow++;
            table.rows[currentRow].cells[currentCol].focus();
        }
        break;
        
        case 32: // clear col
        event.preventDefault();
        table.rows[currentRow].cells[currentCol].innerText=""
        
        break;
        case 13: // enter
        event.preventDefault();


        // 第一格 (產生日期)
        if (currentCol == 0){
            var today = new Date();
            var mmdd = (today.getMonth() + 1).toString().padStart(2, '0') + today.getDate().toString().padStart(2, '0');
            table.rows[currentRow].cells[currentCol].innerText = mmdd;
            currentCol++;
            table.rows[currentRow].cells[currentCol].focus();
            break;
        }

        // 第二格 (魚種)            
        if (currentCol == 1){
            data = table.rows[currentRow].cells[currentCol].innerText
            data=checkFishColor(data)
            
            if (typeof data === 'undefined') {
                break;
            }
            
            table.rows[currentRow].cells[currentCol].innerText=data
            currentCol++;
            table.rows[currentRow].cells[currentCol].focus();
            
            break;
        }

        // 第三格 (重量)            
        if (currentCol == 2){
            data = table.rows[currentRow].cells[currentCol].innerText
            if (data!=""){
            table.rows[currentRow].cells[currentCol].focus();
            currentCol++;
            table.rows[currentRow].cells[currentCol].focus();
            }
            break;
        }

        // 第四格 (單價)            
        if (currentCol == 3){
            data = table.rows[currentRow].cells[currentCol].innerText
            if (data!=""){
            table.rows[currentRow].cells[currentCol].focus();
            currentCol++;
            table.rows[currentRow].cells[currentCol].focus();
            }
            break;
        }

        // 第五格 (分)            
        if (currentCol == 4){

            data = table.rows[currentRow].cells[currentCol].innerText

            if (data==""){
                data = table.rows[currentRow].cells[currentCol].innerText=0
            }

            table.rows[currentRow].cells[currentCol].focus();
            currentCol++;
            table.rows[currentRow].cells[currentCol].focus();
            break;
        }
        // 第六格 (龍)
        if (currentCol == 5){

            if (table.rows[currentRow].cells[currentCol].innerHTML==2){
                table.rows[currentRow].cells[currentCol].innerHTML="小"
            }
            if (table.rows[currentRow].cells[currentCol].innerHTML==3){
                table.rows[currentRow].cells[currentCol].innerHTML="大"
            }
            if (table.rows[currentRow].cells[currentCol].innerHTML=="b"){
                table.rows[currentRow].cells[currentCol].innerHTML="2小"
            }
            if (table.rows[currentRow].cells[currentCol].innerHTML=="c"){
                table.rows[currentRow].cells[currentCol].innerHTML="2大"
            }
            if (table.rows[currentRow].cells[currentCol].innerHTML==""){
                data = table.rows[currentRow].cells[currentCol].innerText=0
            }
            currentCol++;
            table.rows[currentRow].cells[currentCol].focus();
            break;
        }
        // 第七格 (總價)
        if (currentCol == 6){

            // 檢查欄位是否都有數值
            check_empty=1
            for (var i = 0; i < 6;i++){
                if (table.rows[currentRow].cells[i].innerText==""){
                    check_empty=0
                }
            }
            // 所有欄位皆有數值時才進行計算
            if (check_empty ==1){

                // 票據
                ticket=5
                // 單價
                price = table.rows[currentRow].cells[3].innerText
                // 重量
                weight = table.rows[currentRow].cells[2].innerText

                // 小籠 60
                if (table.rows[currentRow].cells[5].innerText == "小"){
                    fish_case=30
                // 大籠 60
                }else if(table.rows[currentRow].cells[5].innerText == "大"){
                    fish_case=60
                // 兩個小籠 30 *2
                }else if(table.rows[currentRow].cells[5].innerText == "2小"){
                    fish_case=60
                // 兩個大籠 60 * 2
                }else if(table.rows[currentRow].cells[5].innerText == "2大"){
                    fish_case=120
                // 如果沒有依照該欄位數值進行加總,預設為 0
                }else{
                    fish_case=table.rows[currentRow].cells[5].innerText
                }

                //  記帳公式 (單價*重量+票據)*1.06  預設票據為五塊錢
                result = (price*weight+ticket)*1.06
                // 將結果加上籠子重量
                result+=parseInt(fish_case);
                // 四捨五入
                result=Math.round(result * 10)/ 10
                // 回填結果
                table.rows[currentRow].cells[currentCol].innerText=result
                // 取得今天帳目
                today_count = current_count.innerHTML

                // 加總帳目 取小數點四捨五入後一位
                today_count = parseFloat(today_count)
                today_count+=result

                // 回填帳目
                current_count.innerHTML=today_count
                
                // 新增下一列,繼續進行運作
                var newRow = table.insertRow(-1);
                for (var i = 0; i < table.rows[0].cells.length; i++) {
                var cell = newRow.insertCell(i);
                cell.contentEditable = true;
                }
                currentRow++;

                // 移動到下一列的第一格
                table.rows[currentRow].cells[0].focus();
                currentCol=0
            }

            break;
        }


        case 83: // clear col
        event.preventDefault();
        button=document.getElementById('myButton');
        button.click()
        alert(83)
        break;
    }
});