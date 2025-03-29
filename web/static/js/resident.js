// Логика для обработки кнопки "Сохранить изменения"
document.getElementById("saveButton").addEventListener("click", function() {
    let residentId = document.getElementById("residentId").value;
    let username = document.getElementById("username").value;
    let email = document.getElementById("email").value;
    let institute = document.getElementById("institute").value;
    let role = document.getElementById("role").value;
    let password = document.getElementById("password").value;

    fetch(`/update_resident/${residentId}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: residentId, username, email, institute, role, password })
    })
    .then(response => {
        if (!response.ok) throw new Error("Ошибка при обновлении данных");
        return response.json();
    })
    .then(data => {
        console.log("Ответ сервера:", data);
        showToast("Данные успешно обновлены!", "success");
    })
    .catch(error => {
        console.error("Ошибка:", error);
        showToast("Не удалось сохранить изменения", "danger");
    });
});

// Функция для показа уведомлений
function showToast(message, type) {
    let toastContainer = document.getElementById("toastContainer");

    let toast = document.createElement("div");
    toast.className = `toast align-items-center text-bg-${type} border-0 show`;
    toast.role = "alert";
    toast.innerHTML = `
        <div class="d-flex">
            <div class="toast-body">${message}</div>
            <button type="button" class="btn-close me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
    `;

    toastContainer.appendChild(toast);

    setTimeout(() => {
        toast.classList.remove("show");
        setTimeout(() => toast.remove(), 500);
    }, 3000);
}       

// Логика для переключения между скрытым и видимым паролем
document.getElementById("togglePassword").addEventListener("click", function() {
    const passwordField = document.getElementById("password");
    const passwordButton = document.getElementById("togglePassword");

    if (passwordField.type === "password") {
        passwordField.type = "text";
        passwordButton.innerHTML = '<i class="bi bi-eye-slash"></i>';
    } else {
        passwordField.type = "password";
        passwordButton.innerHTML = '<i class="bi bi-eye"></i>';
    }
});