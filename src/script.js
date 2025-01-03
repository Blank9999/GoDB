let employeeData = [];

function renderTable() {
  const tableBody = document.querySelector("#employeeTable tbody");
  tableBody.innerHTML = "";

  employeeData.forEach((employee, index) => {
    const row = document.createElement("tr");

    row.innerHTML = `
    <td>${index + 1}</td>
    <td><input name="name" type="text" value="${
      employee.name
    }" data-index="${index}" onchange="updateEmployee(${index}, 'name', this.value)"/></td>

    <td><input name="age" type="text" value="${
      employee.age
    }" data-index="${index}" onchange="updateEmployee(${index}, 'age', this.value)"/></td>

    <td><input name="contact" type="text" value="${
      employee.contact
    }" data-index="${index}" onchange="updateEmployee(${index}, 'contact', this.value)"/></td>

    <td><input name="city" type="text" value="${
      employee.city
    }" data-index="${index}" onchange="updateEmployee(${index}, 'city', this.value)"/></td>

    <td><input name="province" type="text" value="${
      employee.province
    }" data-index="${index}" onchange="updateEmployee(${index}, 'province', this.value)"/></td>

    <td><input name="country" type="text" value="${
      employee.country
    }" data-index="${index}" onchange="updateEmployee(${index}, 'country', this.value)"/></td>

    <td><input name="postalCode" type="text" value="${
      employee.postalCode
    }" onchange="updateEmployee(${index}, 'postalCode', this.value)"/></td>`;

    tableBody.appendChild(row);
  });

  console.log(employeeData);
}

function updateEmployee(index, field, value) {
  employeeData[index][field] = value;
  console.log(employeeData);
}

function addRow() {
  let newEmployee = {
    name: "",
    age: "",
    contact: "",
    city: "",
    province: "",
    country: "",
    postalCode: "",
  };

  employeeData.push(newEmployee);
  renderTable();
}

function onSubmit() {
  //   const data = employeeData;
  const collectionName = document.getElementById("tableName").value;

  const data = {
    collectionName: collectionName, // Include the collection name in the data
    employees: employeeData, // Include the employee data
  };

  console.log("The data is ", data);

  fetch("http://localhost:8080/api/employees", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((res) => res.json())
    .catch((error) => {
      console.log("The data is ", data);
      console.error("Error:", error);
      //   alert("There was an error sending the data.");
    });
}
