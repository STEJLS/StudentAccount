function disp(divname){
    $("#div1").addClass("d-none");
    $("#div2").addClass("d-none");
    $("#div3").addClass("d-none");
    $("#div4").addClass("d-none");
    $("#div5").addClass("d-none");
    $("#div6").addClass("d-none");
    $("#div7").addClass("d-none");
    $("#div8").addClass("d-none");
    $("#chPas").addClass("d-none");
    $("#"+divname).removeClass("d-none");
}

function checkAuth(json) {
    if (json.Message == "Необходимо авторизоваться"){
        window.location.href = "index.html";
    }
}

// выход из ака
function logout() {
fetch('account/logout',{method:"post"})
    .then(r => {
        if (!r.ok){
            console.log("вы вышли");
        }
        r.json().then(json =>{
            checkAuth(json);
            deleteToken()
            window.location.href = "index.html";
        })
    })
    .catch(e => {
        console.log("Произошла ошибка на сервере.");
    })
}
  
function setCookie(name, value, days, path) {    
    path = path || '/'; // заполняем путь если не заполнен
    days = days || 10;  // заполняем время жизни если не получен параметр

    var last_date = new Date();
    last_date.setDate(last_date.getDate() + days);
    var value = escape(value) + ((days==null) ? "" : "; expires="+last_date.toUTCString());
    document.cookie = name + "=" + value + "; path=" + path; // вешаем куки
}

function deleteToken() {
    setCookie('token','11', -20, '/');
}

// смена пароля
function setErrorNote(blockName, errorText) {
$("#"+blockName).text(errorText).removeClass("d-none").removeClass("text-success").addClass("text-danger");      
}

function setSuccessNote(blockName, successText) {
$("#"+blockName).text(successText ).removeClass("d-none").removeClass("text-danger").addClass("text-success");      
}

function clearChangePassForm() {
document.querySelector("#spinner").classList.remove("d-none");
}

function getPasOptions() {
    var form = new FormData();
    form.append("password", $("#oldPassword").val());
    form.append("newPassword", $("#newPassword").val());

    return  {method:"post", body: form }
}

function changePassword(e) {
    e.preventDefault();
    setErrorNote("ChangePassResponse", "");
    if ($("#newPassword").val() != $("#reNewPassword").val() ){
        $("#newPassword").val("");
        $("#reNewPassword").val("");
        setErrorNote("ChangePassResponse", "Указанные новые пароли разные");
        return
    }

fetch('account/changePassword', getPasOptions()).then(r => {
    r.json().then( json =>{
    checkAuth(json);
    if (json.Сompleted){
        setSuccessNote("ChangePassResponse", json.Message);
    } else{
        setErrorNote("ChangePassResponse", json.Message);
        clearChangePassForm();
    }      
    })
}).catch(error => {
        setErrorNote("Ошибка на сервере, попробуйте позже.");   
});    
}

//  Отправка файла 
function getFileOption(e){
const fileInput = e.target.querySelector('#file');
const formData = new FormData();
formData.append('csvFile', fileInput.files[0]);

return options = {
    method: 'POST',
    body: formData,
};
}   

function csvHandler(e, url) {
e.preventDefault();
e.target.querySelector("#response").classList.remove("text-success");
e.target.querySelector("#response").classList.remove("text-danger");
e.target.querySelector('#spinner').classList.remove("d-none");

fetch(url, getFileOption(e)).then( r => {
    r.json().then( json => {
    checkAuth(json);
        e.target.querySelector("#response").innerHTML = json.Message;
        e.target.querySelector("#response").classList.remove("d-none");
    if (json.Сompleted){
        e.target.querySelector("#response").classList.add("text-success");
        e.target.querySelector("#response").classList.remove("text-danger");
    } else{
        e.target.querySelector("#response").classList.add("text-danger");
        e.target.querySelector("#response").classList.remove("text-success");
    }
    e.target.querySelector('#spinner').classList.add('d-none');
    })
}).catch(e => {
    e.target.querySelector("#response").innerHTML = "Ошибка на сервере, попробуйте позже."
    e.target.querySelector("#response").classList.add("text-danger");
    e.target.querySelector("#response").classList.remove("text-success"); 
    e.target.querySelector('#spinner').classList.add('d-none');
}) 
}

function getTempPasswords(e) {

fetch("admin/tempPasswords", {method: 'GET'})
    .then(r => {
        r.blob()
        .then(blob => {
            var url = window.URL.createObjectURL(blob);
            var a = document.createElement('a');
            a.href = url;
            a.download = "passwords.csv";
            document.body.appendChild(a);
            a.click();    
            a.remove();  
        })
    }
    ).catch(e => {

}) 
}


// верификатор
  
  function clearForm() {
      $("#login").val("");
      $("#password").val("");
      $("#fio").val("");
      $("#faculties").val("");
      $("#departments").val("");
  }
  
  function getOptions() {
    var form = new FormData();
    form.append("login", $("#login").val());
    form.append("password", $("#password").val());
    form.append("fullName", $("#fio").val());
    form.append("id_faculty", $("#faculties").val());
    form.append("id_department", $("#departments").val());
    return  {method:"post", body: form }
  }
  
  function test() {
      console.log($("#faculties").val())
      console.log($("#departments").val())
  }
  
  function createVerif(e) {
    e.preventDefault();
    setErrorNote("addVerifResponse", "");
    e.target.querySelector('#spinner').classList.remove("d-none");
  
    if ($("#rePassword").val() != $("#password").val() ){
      $("#rePassword").val("");
      $("#password").val("");
      setErrorNote("addVerifResponse", "Указанные новые пароли разные");
      e.target.querySelector('#spinner').classList.add('d-none');
      return
    }
  
    fetch("admin/verif", getOptions()).then(r => r.json().then( json =>{
      checkAuth(json)
  
      if (!json.Сompleted){
          setErrorNote("addVerifResponse", json.Message);
      }else{
          setSuccessNote("addVerifResponse", json.Message);
          clearForm();
      }
      e.target.querySelector('#spinner').classList.add('d-none');
    }))
    .catch(e =>{
      setErrorNote("addVerifResponse", "Ошибка на сервере. Повторите позже.");
      e.target.querySelector('#spinner').classList.add('d-none');
    })
  }
  
  
  
  function getDepartments(event) {
  
      fucID = event.target.value;
      updateDepartments(event.target.value);
  }
  
  function updateDepartments(idFuc) {
      setErrorNote("addVerifResponse", "");
  
      fetch('admin/departments').then(r => {
        r.json().then(json => {
          checkAuth(json)
            if (!json.Сompleted){
              setErrorNote("addVerifResponse", json.Message);
              return;
            }
          const selectDep = document.getElementById("departments");
  
          while (selectDep.firstChild) {
              selectDep.removeChild(selectDep.firstChild);
          }        
  
          json.Body.forEach(item => {
              if (item.IDFaculty != idFuc){
                  return;
              }
                  
  
              var newOption = document.createElement('option');
              newOption.innerHTML = item.Name;
              newOption.value = item.ID;
              selectDep.appendChild(newOption);
          })
              
        })
      }).catch(error => {
          setErrorNote("Ошибка на сервере, попробуйте позже.");   
    });  
  }


  function FOSandRPDHandler(e, url) {
    e.preventDefault();
    e.target.querySelector("#response").classList.remove("text-success");
    e.target.querySelector("#response").classList.remove("text-danger");
    e.target.querySelector('#spinner').classList.remove("d-none");
    
    fetch(url).then( r => {
        r.json().then( json => {
        checkAuth(json);
            e.target.querySelector("#response").innerHTML = json.Message;
            e.target.querySelector("#response").classList.remove("d-none");
        if (json.Сompleted){
            e.target.querySelector("#response").classList.add("text-success");
            e.target.querySelector("#response").classList.remove("text-danger");
        } else{
            e.target.querySelector("#response").classList.add("text-danger");
            e.target.querySelector("#response").classList.remove("text-success");
        }
        e.target.querySelector('#spinner').classList.add('d-none');
        })
    }).catch(e => {
        e.target.querySelector("#response").innerHTML = "Ошибка на сервере, попробуйте позже."
        e.target.querySelector("#response").classList.add("text-danger");
        e.target.querySelector("#response").classList.remove("text-success"); 
        e.target.querySelector('#spinner').classList.add('d-none');
    }) 
    }