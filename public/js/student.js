function disp(divname){
    $("#div1").addClass("d-none");
    $("#div2").addClass("d-none");
    $("#div3").addClass("d-none");
    $("#div4").addClass("d-none");
    $("#div5").addClass("d-none");
    $("#div6").addClass("d-none");
    $("#div7").addClass("d-none");
    $("#chPas").addClass("d-none");
    $("#"+divname).removeClass("d-none");
}

function logout() {
fetch('account/logout',{method:"post"})
    .then(r => r.json().then(json =>{
            checkAuth(json);
            deleteToken()
            window.location.href = "index.html";
        })
    )
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

function checkAuth(json) {
    if (json.Message == "Необходимо авторизоваться"){
        window.location.href = "index.html";
    }
}

function setErrorNote(blockName, errorText) {
    $("#"+blockName).text(errorText).removeClass("d-none").removeClass("text-success").addClass("text-danger");      
}

function setSuccessNote(blockName, successText) {
    $("#"+blockName).text(successText ).removeClass("d-none").removeClass("text-danger").addClass("text-success");      
}

function getArticleOptions() {
    const fileInput = document.querySelector('#articleFile');

    var form = new FormData();
    form.append("name", $("#articleName").val());
    form.append("journal", $("#journalName").val());
    form.append("biblioRecord", $("#biblioRecord").val());
    form.append("type", $("#articleType").val());
    form.append("article", fileInput.files[0]);
    return  {method:"post", body: form }
}

function createArticle(e) {
    e.preventDefault();    
    setErrorNote("createArticleResponse", "");

      fetch('student/article', getArticleOptions()).then( r => {
        r.json().then(json => {
            checkAuth(json);
            if (json.Сompleted){
                setSuccessNote("createArticleResponse", json.Message);
                e.target.reset();
            } else{
                setErrorNote("createArticleResponse", json.Message);
            }      
        }).catch(error => {
            setErrorNote("createArticleResponse", "Ошибка на сервере, попробуйте позже.");   
        });   
      })
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
            setErrorNote("ChangePassResponse","Ошибка на сервере, попробуйте позже.");   
    });    
}

function setPersonalInfo() {
    fetch('student/info').then(r => {
        r.json().then( json =>{
            checkAuth(json);
            var s = json.Body;
            $("#fio").text(s.FullName+"("+s.Number+")");
            $("#faculty").text(s.FacultyName+"("+s.FacultyShortName+")");
            $("#department").text(s.DepartmentName+"("+s.DepartmentShortName+")");
            $("#fieldOfStudy").text(s.FieldName+"("+s.FieldCode+")");
            $("#group").text(s.Team+"-"+s.TeamNumber);
            $("#brand").text("КФ МГТУ"+" "+s.FullName);
    })});
}

function setMarks() {
    fetch('student/marks').then(r => {
        r.json().then( json =>{
            checkAuth(json);
            var s = json.Body;

            const table = document.getElementById("progress");
            json.Body.forEach(item => {
                var newtr = document.createElement('tr');
                var subject = document.createElement('td');
                var mark = document.createElement('td');
                var type = document.createElement('td');
                var re = document.createElement('td');
                var semester = document.createElement('td');
                subject.innerHTML = item.Subject;
                mark.innerHTML = item.Rating;

                var passType;
                switch (item.PassType) {
                    case 0:
                    passType = "Экзамен"
                        break;
                    case 1:
                    passType = "Д.зачет"
                        break;
                    case 2:
                    passType = "Зачет"
                        break;
                    
                    default:
                        break;
                }
                type.innerHTML = passType;

                if(item.PassType == 2){
                    if (item.Rating == 0){
                        mark.innerHTML = "Не сдан";
                    }else{
                        mark.innerHTML = "Cдан";
                    }
                }
                re.innerHTML = item.Repass? "Да" : "Нет";
                semester.innerHTML = item.Semester;

                newtr.appendChild(subject);
                newtr.appendChild(mark);
                newtr.appendChild(type);
                newtr.appendChild(re);
                newtr.appendChild(semester);
                table.appendChild(newtr);
            })
    })});
}

function setPractices() {
    fetch('/student/practices').then(r => {
        r.json().then( json =>{
            checkAuth(json);
            var a = json.Body;

            const allPractices = document.getElementById("div3");
            var count = 1;

            json.Body.forEach(item => {
                var newp = document.createElement('p');
                newp.classList.add("h5", "mt-4", "text-primary" );
                newp.innerHTML = "Семетр " + item.Semester; 
                allPractices.appendChild(newp);

                var newTable = document.createElement('table');
                newTable.classList.add("table", "table-bordered", "table-sm");
                
                var newTbody = document.createElement('tbody');
                
                var newTr = document.createElement('tr');
                var newTh = document.createElement('th');
                var newTd = document.createElement('td');

                newTh.innerHTML = "Название";
                newTh.style = "width:20%";
                newTd.innerHTML = item.Name;

                newTr.appendChild(newTh);
                newTr.appendChild(newTd);                

                newTbody.appendChild(newTr);

                newTr = document.createElement('tr');
                newTh = document.createElement('th');
                newTd = document.createElement('td');

                newTh.innerHTML = "Руководитель";
                newTd.innerHTML = item.Head;

                newTr.appendChild(newTh);
                newTr.appendChild(newTd);                

                newTbody.appendChild(newTr);

                newTr = document.createElement('tr');
                newTh = document.createElement('th');
                newTd = document.createElement('td');

                newTh.innerHTML = "Предпириятие";
                newTd.innerHTML = item.Company;

                newTr.appendChild(newTh);
                newTr.appendChild(newTd);                

                newTbody.appendChild(newTr);
                
                newTr = document.createElement('tr');
                newTh = document.createElement('th');
                newTd = document.createElement('td');

                newTh.innerHTML = "Дата";
                newTd.innerHTML = item.Date;

                newTr.appendChild(newTh);
                newTr.appendChild(newTd);                

                newTbody.appendChild(newTr);
                
                newTr = document.createElement('tr');
                newTh = document.createElement('th');
                newTd = document.createElement('td');

                newTh.innerHTML = "Оценка";
                newTd.innerHTML = item.Rating;

                newTr.appendChild(newTh);
                newTr.appendChild(newTd);                

                newTbody.appendChild(newTr);

                newTable.appendChild(newTbody);
                allPractices.appendChild(newTable);
     
                count++;
            })

            if (count == 1){
                $("#noPracticesToShow").removeClass("d-none");
            }
    })});
}

function setArticles() {
    fetch('/student/articles').then(r => {
        r.json().then( json =>{
            checkAuth(json);
            var a = json.Body;

            const confirmedTable = document.getElementById("confirmedArticleTable");
            const notConfirmedTable = document.getElementById("notConfirmedArticleTable");
            var confCount = 1;
            var noConfCount = 1;
            
            while (confirmedTable.firstChild) {
                confirmedTable.removeChild(confirmedTable.firstChild);
            }
            while (notConfirmedTable.firstChild) {
                notConfirmedTable.removeChild(notConfirmedTable.firstChild);
            }
            
            json.Body.forEach(item => {
                var newtr = document.createElement('tr');
                var newth = document.createElement('th');
                var newtd = document.createElement('td');
                var newa = document.createElement('a');
                newth.style = "width:5%";
                newa.innerHTML = item.Name;
                newa.setAttribute("href", "");
                newa.setAttribute("data-toggle", "modal");
                newa.setAttribute("data-target", "#article-modal");
                newa.setAttribute("data-name", item.Name);
                newa.setAttribute("data-journal", item.Journal);
                newa.setAttribute("data-ref", item.BiblioRecord);
                newa.setAttribute("data-type", item.ArticlType);
                newa.setAttribute("data-id", item.ID);
                newtd.appendChild(newa);
                newtr.appendChild(newth);
                newtr.appendChild(newtd);

                if (item.Confirmed){
                    $("#confirmedArticleTitle").removeClass("d-none");
                    newth.innerHTML = confCount++;
                    confirmedTable.appendChild(newtr);      
                }else{
                    $("#notConfirmedArticleTitle").removeClass("d-none");
                    newth.innerHTML = noConfCount++;
                    notConfirmedTable.appendChild(newtr);                     
                }        
            })

            if (confCount == 1 && noConfCount == 1){
                $("#noArticlesToShow").removeClass("d-none");
            }
    })});
}

function setHandlerForModalArticle(event) {
    console.log("fd");
    var button = $(event.relatedTarget) // Кнопка, что спровоцировало модальное окно  
    var name = button.data('name') 
    var journal = button.data('journal') 
    var ref = button.data('ref') 
    var type = button.data('type') 
    var id = button.data('id') 

    var modal = $(this)
    modal.find('#article-name').val(name)
    modal.find('#article-journal').val(journal)
    modal.find('#article-ref').val(ref)
    modal.find('#article-type').val(type)
    modal.find('#article-download').attr("href", "student/article/"+id)
}