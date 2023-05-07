const table = document.getElementById("customersTable");
const url = "/get_customer_name";

fetch(url)
  .then(response => response.json())
  .then(data => {
    let row = table.insertRow();
    let count = 0;
    
    data.forEach(customer => {
      if (count === 6) {
        row = table.insertRow();
        count = 0;
      }
      const cell = row.insertCell();
      cell.innerHTML = customer.name;
      count++;
    });
    
    // If the last row isn't full, add empty cells to it
    while (count < 6) {
      const cell = row.insertCell();
      cell.innerHTML = "&nbsp;";
      count++;
    }
  });


var currentRow = 1;
var currentCol = 0;


var table2 = document.getElementById("customersTable");

table2.addEventListener("keydown", function(event) {

table2.rows[1].cells[0].focus();


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
}
});