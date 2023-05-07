let dictionary = {};

window.onload = function() {
    var table = document.getElementById("myTable");
    table.rows[1].cells[0].focus();

    // 宣告 dictionary 變數
    var customer = document.getElementById("customer");
    customer.innerHTML+="測試客戶";

    var current_count = document.getElementById("current_count");
    current_count.innerHTML+=100;

    var pre_count = document.getElementById("pre_count");
    pre_count.innerHTML+=100;



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
        
            if(checkFishColor(data)!="未知"){
            table.rows[currentRow].cells[currentCol].innerText=checkFishColor(data)
            currentCol++;
            table.rows[currentRow].cells[currentCol].focus();
            break;
            }

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
            if (data==""){
                data = table.rows[currentRow].cells[currentCol].innerText=0
            }
            currentCol++;
            table.rows[currentRow].cells[currentCol].focus();
            break;
        }
        // 第七格 (總價)
        if (currentCol == 6){
            // 票據
            ticket=5
            price = table.rows[currentRow].cells[3].innerText
            weight = table.rows[currentRow].cells[2].innerText

            if (table.rows[currentRow].cells[5].innerText == "小"){
                fish_case=30
            }else if(table.rows[currentRow].cells[5].innerText == "大"){
                fish_case=60
            }else if(table.rows[currentRow].cells[5].innerText == "2小"){
                fish_case=60
            }else if(table.rows[currentRow].cells[5].innerText == "2大"){
                fish_case=120
            }else{
                fish_case=table.rows[currentRow].cells[5].innerText
            }

            result = (price*weight+ticket)*1.06
            result+=parseInt(fish_case);
            result=Math.round(result * 10)/ 10
            table.rows[currentRow].cells[currentCol].innerText=result

            // 當第一行欄位沒有數值時,不再移動到下一行
            if (table.rows[currentRow].cells[0].innerText!=""){
                var newRow = table.insertRow(-1);
                for (var i = 0; i < table.rows[0].cells.length; i++) {
                var cell = newRow.insertCell(i);
                cell.contentEditable = true;
                }
                currentRow++;
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