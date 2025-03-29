function filterStatements(button) {
    // Убираем активный класс у всех кнопок
    document.querySelectorAll('.btn-group .btn').forEach(btn => {
        btn.classList.remove('active');
    });
    
    // Добавляем активный класс нажатой кнопке
    button.classList.add('active');
    
    const filterStatus = button.getAttribute('data-status').toLowerCase();
    const statements = document.querySelectorAll(".statement-item");
    
    statements.forEach(function(statement) {
        const status = statement.getAttribute("data-status").toLowerCase();
        if (filterStatus === "" || status.includes(filterStatus)) {
            statement.style.display = "";
        } else {
            statement.style.display = "none";
        }
    });
}

// Добавляем обработчик события при загрузке страницы
document.addEventListener("DOMContentLoaded", function() {
    // Активируем фильтр "Все" при загрузке страницы
    const filterButton = document.querySelector('.btn-group .btn[data-status=""]');
    if (filterButton) {
        filterStatements(filterButton);
    }
});