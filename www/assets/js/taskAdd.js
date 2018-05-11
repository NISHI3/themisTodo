var taskAddForm = document.querySelector("#taskboardAddForm"),
    taskAddFormErrorElem = taskAddForm.querySelector(".error"),
    taskAddShow = document.querySelector("#taskboardAdd"),
    closeFormElem = document.querySelector("#taskboardAddClose");

taskAddForm.addEventListener("submit", postTaskAdd, true);
taskAddFormErrorElem.addEventListener("click", clickError, true);
taskAddShow.addEventListener("click", taskAddShowClick, true);
closeFormElem.addEventListener("click", clickCloseForm, true);

var nowTime = new Date(Date.now());
taskAddForm.querySelector("input[name=deadline]").value = dateFormat(nowTime);


function postTaskAdd(e) {
    e.preventDefault();

    let formData = new FormData(taskAddForm);
    let projectAddJson = {
        "name": formData.get("name"),
        "deadline": formData.get("deadline"),
        "description": formData.get("description"),
        "projectId": projectId
    };

    TaskApi.Create(projectAddJson).then(function (json) {
        if (!json.success) {
            taskAddFormErrorElem.style.display = "block";
            taskAddFormErrorElem.innerText = json.message;
        } else {
            taskAddFormErrorElem.style.display = "none";
            addFormClear();
            TaskApi.GetTaskFromCreateDate(json.createDate).then(function (json) {

            });
        }
    });
}

function clickError(e) {
    e.target.style.display = "none";
}

function taskAddShowClick(e) {
    e.preventDefault();
    taskAddForm.style.right = "0";
}

function clickCloseForm() {
    taskAddForm.style.right = "-" + taskAddForm.clientWidth + "px";
}

function addFormClear(){
    taskAddForm.reset();
    var nowTime = new Date(Date.now());
    taskAddForm.querySelector("input[name=deadline]").value = dateFormat(nowTime);
}

function dateFormat(date) {
    var y = date.getFullYear();
    var m = date.getMonth() + 1;
    var d = date.getDate();

    if (m < 10) {
        m = '0' + m;
    }
    if (d < 10) {
        d = '0' + d;
    }

    return y + '-' + m + '-' + d;
}