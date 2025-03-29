document.addEventListener("DOMContentLoaded", function () {
    let form = document.getElementById("serviceForm");
    let saveButton = document.getElementById("saveButton");
    
    if (form) {
        form.addEventListener("input", function () {
            saveButton.style.display = "block";
        });
    }
});