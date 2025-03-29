function applyFilters() {
    let statusFilter = document.getElementById("filterStatus").value.toLowerCase();
    let typeFilter = document.getElementById("filterType").value.toLowerCase();
    let hostelFilter = document.getElementById("filterHostel").value.toLowerCase();
    let searchRoomNumber = document.getElementById("searchRoomNumber").value.toLowerCase();

    let cards = document.querySelectorAll(".room-card");

    cards.forEach(card => {
        let status = card.querySelector(".badge").textContent.toLowerCase();
        let type = card.querySelector(".list-group-item:nth-child(1)").textContent.toLowerCase();
        let hostel = card.querySelector(".list-group-item:nth-child(4)").textContent.toLowerCase();
        let roomNumber = card.querySelector(".card-title").textContent.toLowerCase();

        if ((statusFilter === "" || status.includes(statusFilter)) &&
            (typeFilter === "" || type.includes(typeFilter)) &&
            (hostelFilter === "" || hostel.includes(hostel)) &&
            (searchRoomNumber === "" || roomNumber.includes(searchRoomNumber))) {
            card.style.display = "";
        } else {
            card.style.display = "none";
        }
    });
}