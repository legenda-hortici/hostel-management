function applyFilters() {
    let statusFilter = document.getElementById("filterStatus").value.toLowerCase();
    let typeFilter = document.getElementById("filterType").value.toLowerCase();
    let hostelFilter = document.getElementById("filterHostel") ? document.getElementById("filterHostel").value.toLowerCase() : "";
    let searchRoomNumber = document.getElementById("searchRoomNumber").value.toLowerCase();

    let cards = document.querySelectorAll(".room-card");

    cards.forEach(card => {
        // Получаем текст из карточки
        let roomNumber = card.querySelector(".card-title").textContent.toLowerCase().replace("комната №", "").trim();

        let typeItem = card.querySelector(".list-group-item:nth-child(1)").textContent.toLowerCase().replace("тип:", "").trim();
        let statusItem = card.querySelector(".list-group-item:nth-child(2) .badge").textContent.toLowerCase().trim();
        let hostelItem = card.querySelector(".list-group-item:nth-child(3)").textContent.toLowerCase().replace("общежитие:", "").trim();

        // Проверка соответствия фильтрам
        let match =
            (statusFilter === "" || statusItem.includes(statusFilter)) &&
            (typeFilter === "" || typeItem.includes(typeFilter)) &&
            (hostelFilter === "" || hostelItem.includes(hostelFilter)) &&
            (searchRoomNumber === "" || roomNumber.includes(searchRoomNumber));

        card.style.display = match ? "" : "none";
    });
}
