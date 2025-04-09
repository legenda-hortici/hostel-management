document.getElementById("saveButton").addEventListener("click", function() {
    const residentId = document.getElementById("residentId").value;
    
    const formData = {
        username: document.getElementById("username").value,
        surname: document.getElementById("surname").value,
        email: document.getElementById("email").value,
        institute: document.getElementById("institute").value,
        role: document.getElementById("role").value,
        password: document.getElementById("password").value
    };

    fetch(`/admin/residents/resident/${residentId}/edit`, {
        method: 'PUT',
        headers: { 
            'Content-Type': 'application/json',
            'X-HTTP-Method-Override': 'PUT'
        },
        body: JSON.stringify(formData)
    })
    .then(async response => {
        const text = await response.text();
        
        try {
            const data = text ? JSON.parse(text) : {};
            
            if (!response.ok) {
                const error = new Error(data.message || 'Ошибка сервера');
                error.response = data;
                throw error;
            }
            
            return data;
        } catch (e) {
            if (e instanceof SyntaxError) {
                throw new Error(`Невалидный JSON ответ: ${text.substring(0, 100)}...`);
            }
            throw e;
        }
    })
    .then(data => {
        console.log("Успешный ответ:", data);
        showToast("Данные успешно обновлены!", "success");
        setTimeout(() => window.location.reload(), 1500);
    })
    .catch(error => {
        console.error("Полная ошибка:", {
            message: error.message,
            stack: error.stack,
            response: error.response
        });
        
        const errorMsg = error.message || "Не удалось сохранить изменения";
        showToast(errorMsg, "danger");
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