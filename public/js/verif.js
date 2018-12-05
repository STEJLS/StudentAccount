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
                e.target.reset();
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
    $("#brand").text("КФ МГТУ"+" "+JSON.parse(localStorage.getItem('user')).FullName);
}

function setArticlesToVerif() {
    fetch('/verif/articlesForVerif').then(r => {
        r.json().then( json =>{
            checkAuth(json);
            var a = json.Body;

            const table = document.getElementById("articles-to-verif-table");
            var count = 1;
            
            while (table.firstChild) {
                table.removeChild(table.firstChild);
            }
            
            json.Body.forEach(item => {
                var newtr = document.createElement('tr');
                var newth = document.createElement('th');
                var newtd = document.createElement('td');
                var newa = document.createElement('a');
                newth.innerHTML = count++;
                newa.innerHTML = item.Team+" " + getShortName(item.FIO)+" "+item.Name;
                newa.setAttribute("href", "");
                newa.setAttribute("data-toggle", "modal");
                newa.setAttribute("data-target", "#verif-article-modal");
                newa.setAttribute("data-fio", item.FIO);
                newa.setAttribute("data-team", item.Team);
                newa.setAttribute("data-name", item.Name);
                newa.setAttribute("data-journal", item.Journal);
                newa.setAttribute("data-ref", item.BiblioRecord);
                newa.setAttribute("data-type", item.ArticlType);
                newa.setAttribute("data-id", item.ID);
                newtd.appendChild(newa);
                newtr.appendChild(newth);
                newtr.appendChild(newtd);
                table.appendChild(newtr);      
                
            })

            if (count == 1){
                $("#noArticlesForVerir").removeClass("d-none");
            }
    })});
}


function getShortName(name) {
    var splitted = name.split(" ")
    if (splitted.length < 3){
        return name;
    } 
    return `${splitted[0]} ${splitted[1][0]}.${splitted[2][0]}.`;
}


function setHandlerForModalVerifArticle(event) {
        var button = $(event.relatedTarget) // Кнопка, что спровоцировало модальное окно  

        var fio = button.data('fio') 
        var team = button.data('team') 
        var name = button.data('name') 
        var journal = button.data('journal') 
        var ref = button.data('ref') 
        var type = button.data('type') 
        var id = button.data('id') 

        var modal = $(this)
        modal.find('#verif-article-author').val(fio)
        modal.find('#verif-article-team').val(team)
        modal.find('#verif-article-name').val(name)
        modal.find('#verif-article-journal').val(journal)
        modal.find('#verif-article-ref').val(ref)
        modal.find('#verif-article-type').val(type)
        modal.find('#verif-article-download').attr("href", "verif/article/"+id)
        modal.find('#verif-article-cancel-btn').attr("onclick", "verifArticleCancel("+id+")")
        modal.find('#verif-article-confirm-form').attr("onsubmit", "verifArticleConfirm(event,"+id+")")
}

function getCancelArticleOptions(id) {
    var form = new FormData();
    form.append("id", id);

    return  {method:"post", body: form }
}

function verifArticleCancel(id) {
    
    fetch('verif/cancelArticle', getCancelArticleOptions(id)).then(r => {
        r.json().then( json =>{
            checkAuth(json);
            if (json.Сompleted){
                setArticlesToVerif();
                setSuccessNote("verifArticleResponse", json.Message);
                $('#verif-article-modal').modal('hide'); 
        } else{
            setErrorNote("verifArticleResponse", json.Message);
        }
    })
    }).catch(error => {
            setErrorNote("verifArticleResponse","Ошибка на сервере, попробуйте позже.");   
    });   
}



function getConfirmArticleOptions(id) {
    var form = new FormData();
    form.append("name", $("#verif-article-name").val());
    form.append("journal", $("#verif-article-journal").val());
    form.append("biblioRecord", $("#verif-article-ref").val());
    form.append("type", $("#verif-article-type").val());
    form.append("id", id);

    return  {method:"post", body: form }
}

function verifArticleConfirm(e, id) {
    e.preventDefault();   
    fetch('verif/article', getConfirmArticleOptions(id)).then(r => {
        r.json().then( json =>{
            checkAuth(json);
            if (json.Сompleted){
                setArticlesToVerif();
                setSuccessNote("verifArticleResponse", json.Message);
                $('#verif-article-modal').modal('hide'); 
            } else{
                setErrorNote("verifArticleResponse", json.Message);
            }
        })
        }).catch(error => {
                setErrorNote("verifArticleResponse","Ошибка на сервере, попробуйте позже.");   
        });    
}