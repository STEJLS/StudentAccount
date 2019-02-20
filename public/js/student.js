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
    .catch(() => {
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
                setArticles();
            } else{
                setErrorNote("createArticleResponse", json.Message);
            }      
        }).catch(() => {
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
    }).catch(() => {
            setErrorNote("ChangePassResponse", "Ошибка на сервере, попробуйте позже.");
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
            $("#profile").text(s.FieldProfile);
            $("#group").text(s.Team+"-"+s.TeamNumber);
            $("#brand").text("КФ МГТУ"+" "+s.FullName);
    })});
}

function setMarks() {
    fetch('student/marks').then(r => {
        r.json().then( json =>{
            checkAuth(json);

            json.Body.sort(function(obj1, obj2) {
                if (obj1.Semester < obj2.Semester) return -1;
                if (obj1.Semester > obj2.Semester) return 1;
                return 0;
            })

            var semesters = [...new Set(json.Body.map(function(item) {
                return item.Semester;
              }))];

            var block = document.getElementById('learningAchievement');  

            for (var i = 0; i < semesters.length; i++) {
                var newTable = `
                <p class="h4">`+(i+1)+` семестр </p>
                <table class="table mt-4">
                  <thead>
                    <tr>
                      <th>Предмет</th>
                      <th>Оценка</th>
                      <th>Тип сдачи</th>
                      <th>Пересдача</th>
                    </tr>
                  </thead>
                  <tbody id="learningAchievementtable`+(i+1)+`"> 
                  </tbody>
                </table>`;
                block.innerHTML += newTable;
            }



            const table = document.getElementById("progress");
            json.Body.forEach(item => {
                var newtr = document.createElement('tr');
                var subject = document.createElement('td');
                var mark = document.createElement('td');
                var type = document.createElement('td');
                var re = document.createElement('td');                
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

                newtr.appendChild(subject);
                newtr.appendChild(mark);
                newtr.appendChild(type);
                newtr.appendChild(re);
                document.getElementById("learningAchievementtable"+item.Semester).appendChild(newtr);
            })
    })});
}

function setPractices() {
    fetch('/student/practices').then(r => {
        r.json().then( json =>{
            checkAuth(json);

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
                newa.setAttribute("data-fileName", item.FileName);
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

function setCourse() {
    fetch('/student/courses').then(r => {
        r.json().then( json =>{
            checkAuth(json);

            const setThemeCourseTable = document.getElementById("setThemeCourseTable");
            const courseTable = document.getElementById("courseTable");
            const notConfirmedCourseTable = document.getElementById("notConfirmedCourseTable");
            var setThemeCount = 1;
            var courseCount = 1;
            var noConfCount = 1;
            
            while (setThemeCourseTable.firstChild) {
                setThemeCourseTable.removeChild(setThemeCourseTable.firstChild);
            }
            while (courseTable.firstChild) {
                courseTable.removeChild(courseTable.firstChild);
            }
            while (notConfirmedCourseTable.firstChild) {
                notConfirmedCourseTable.removeChild(notConfirmedCourseTable.firstChild);
            }
            
            json.Body.forEach(item => {
                var newtr = document.createElement('tr');
                var newth = document.createElement('th');
                var newtd = document.createElement('td');
                var newa = document.createElement('a');
                newth.style = "width:5%";
                newa.innerHTML = item.Subject;
                newa.setAttribute("href", "");
                newa.setAttribute("data-toggle", "modal");
                newa.setAttribute("data-target", "#verif-course-modal");
                newa.setAttribute("data-fio", item.FIO);
                newa.setAttribute("data-team", item.Team);
                newa.setAttribute("data-theme", item.Theme.String);
                newa.setAttribute("data-semester", item.Semester);
                newa.setAttribute("data-head", item.Head);
                newa.setAttribute("data-rating", item.Rating);
                newa.setAttribute("data-id", item.ID);
                newa.setAttribute("data-subject", item.Subject);
                newa.setAttribute("data-confirmed", item.Confirmed);
                newtd.appendChild(newa);
                newtr.appendChild(newth);
                newtr.appendChild(newtd);
         

                if (item.Confirmed && item.Theme.String!=""){
                    $("#courseTitle").removeClass("d-none");
                    newth.innerHTML = courseCount++;
                    courseTable.appendChild(newtr);      
                }else if (!item.Confirmed && item.Theme.String!=""){
                    $("#notConfirmedCourseTitle").removeClass("d-none");
                    newth.innerHTML = noConfCount++;
                    notConfirmedCourseTable.appendChild(newtr);                     
                }else{
                    $("#setThemeCourseTitle").removeClass("d-none");
                    newth.innerHTML = setThemeCount++;
                    setThemeCourseTable.appendChild(newtr);                     
                }
            })

            if (courseCount == 1 && noConfCount == 1 && setThemeCount == 1){
                $("#noCoursesToShow").removeClass("d-none");
            }
    })});
}

function setHandlerForModalArticle(event) {
    document.getElementById("article-download").style.visibility = "visible";
    var button = $(event.relatedTarget) // Кнопка, что спровоцировало модальное окно  
    var name = button.data('name') 
    var journal = button.data('journal') 
    var ref = button.data('ref') 
    var type = button.data('type') 
    var id = button.data('id') 
    var fileName = button.data('filename') 

    var modal = $(this)
    modal.find('#article-name').val(name)
    modal.find('#article-journal').val(journal)
    modal.find('#article-ref').val(ref)
    modal.find('#article-type').val(type)
    modal.find('#article-download').attr("href", "student/article/"+id)

    if (fileName == ""){
        document.getElementById("article-download").style.visibility = "hidden";
    }
}

function setHandlerForModalVerifCourse(event) {
    var button = $(event.relatedTarget) // Кнопка, что спровоцировало модальное окно  

    var fio = button.data('fio') 
    var team = button.data('team') 
    var theme = button.data('theme') 
    var semester = button.data('semester') 
    var head = button.data('head') 
    var rating = button.data('rating') 
    var id = button.data('id') 
    var subject = button.data('subject') 
    
    var modal = $(this)

    modal.find('#verif-course-theme').val(theme);
    if(theme != ""){
        document.getElementById("submitCourseBTN").classList.add('d-none');
        document.getElementById("verif-course-ModalLabel").innerHTML = "Курсовая работа";
        modal.find('#verif-course-theme').attr("readonly","readonly");
    }else{        
        $('#verif-course-theme').removeAttr("readonly");
    }

    modal.find('#verifCourseResponse').addClass("d-none");
    modal.find('#verif-course-author').val(fio);
    modal.find('#verif-course-team').val(team);
    modal.find('#verif-course-head').val(head);
    modal.find('#verif-course-semester').val(semester);
    modal.find('#verif-course-rating').val(rating);
    modal.find('#verif-course-subject').val(subject);
    modal.find('#verif-course-confirm-form').attr("onsubmit", "studentCourseConfirm(event,"+id+")");
}

function getConfirmCourseOptions(id) {
    var form = new FormData();
    form.append("theme", $("#verif-course-theme").val());
    form.append("id", id);

    return  {method:"post", body: form }
}

function studentCourseConfirm(e, id) {
    e.preventDefault();   
    fetch('student/courseWork', getConfirmCourseOptions(id)).then(r => {
        r.json().then( json =>{
            checkAuth(json);
            if (json.Сompleted){
                setCourse();
                setSuccessNote("verifCourseResponse", json.Message);
                $('#verif-course-modal').modal('hide'); 
            } else{
                setErrorNote("verifCourseResponse", json.Message);
            }
        })
        }).catch(() => {
                setErrorNote("verifCourseResponse", "Ошибка на сервере, попробуйте позже.");
            }); 
}

function setProgramsOfDisciplines() {
    fetch('student/FOSandRPDList').then(r => {
        r.json().then( json =>{
            checkAuth(json);
            
            var container = document.getElementById("programsOfDiscipline");
            if (json.Body.length == 0){
                var p = document.createElement('p');
                p.innerHTML = json.Message;
                container.appendChild(p);
                return;
            }

            var table = document.createElement('table')
            table.setAttribute("class", "table")
            var thread = document.createElement('thead')
            var tr = document.createElement('tr')
            var th1 = document.createElement('th')
            var th2 = document.createElement('th')
            var th3 = document.createElement('th')
            var th4 = document.createElement('th')
            th1.innerHTML = "#";
            th1.setAttribute("scope", "col")
            tr.appendChild(th1);
            th2.innerHTML = "Предмет";
            th2.setAttribute("scope", "col")
            tr.appendChild(th2);
            th3.innerHTML = "ФОС";
            th3.setAttribute("scope", "col")
            tr.appendChild(th3)
            th4.innerHTML = "РПД";
            th4.setAttribute("scope", "col")
            tr.appendChild(th4);
            thread.appendChild(tr);
            table.appendChild(thread);            

            var tbody = document.createElement('tbody')
            var count = 1;
            json.Body.forEach(item => {
                var newtr = document.createElement('tr');
                var newth = document.createElement('th');
                newth.setAttribute("scope", "row");
                newth.innerHTML = count++;
                var name = document.createElement('td');
                var fos = document.createElement('td');
                var rpd = document.createElement('td');
                var fosa = document.createElement('a');
                var rpda = document.createElement('a');
                name.innerHTML = item.Name;
                
                fosa.innerHTML = "Скачать";
                rpda.innerHTML = "Скачать";
                
                fosa.setAttribute("href", "./student/document/"+item.FosID);
                rpda.setAttribute("href", "./student/document/"+item.RpdID);
                fos.appendChild(fosa)
                rpd.appendChild(rpda)

                newtr.appendChild(newth)
                newtr.appendChild(name)
                newtr.appendChild(fos)
                newtr.appendChild(rpd)   
                
                tbody.appendChild(newtr);
                
            })
            table.appendChild(tbody)
            container.appendChild(table)
    })});
}