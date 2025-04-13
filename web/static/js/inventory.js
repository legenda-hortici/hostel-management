function applyFilters() {
    const furnitureSelect = document.getElementById("filterInstitute");
    const hostelSelect = document.getElementById("filterHostel");
    const rows = document.querySelectorAll("#residentsTable tbody tr");

    const selectedFurniture = furnitureSelect.value.toLowerCase();
    const selectedHostel = hostelSelect.value;

    rows.forEach(row => {
        const furnitureType = row.querySelector(".inv-name").textContent.toLowerCase().trim();
        const hostelNumber = row.querySelector("td:nth-child(4)").textContent.trim();

        const matchesFurniture = selectedFurniture === "" || furnitureType === selectedFurniture;
        const matchesHostel = selectedHostel === "" || hostelNumber === selectedHostel;

        if (matchesFurniture && matchesHostel) {
            row.style.display = "";
        } else {
            row.style.display = "none";
        }
    });
}

document.getElementById("searchInput").addEventListener("input", function() {
    const searchTerm = this.value.toLowerCase();
    const rows = document.querySelectorAll("#residentsTable tbody tr");

    rows.forEach(row => {
        const roomNumber = row.querySelector(".room-number").textContent.toLowerCase();
        const invNumber = row.querySelector("td:nth-child(4)").textContent.toLowerCase();

        if (roomNumber.includes(searchTerm) || invNumber.includes(searchTerm)) {
            row.style.display = "";
        } else {
            row.style.display = "none";
        }
    });
});

document.addEventListener('DOMContentLoaded', function () {
    const modal = document.getElementById('editInventoryModal');
    modal.addEventListener('show.bs.modal', function (event) {
        const button = event.relatedTarget;

        document.getElementById('modal-id').value = button.getAttribute('data-inv-id');
        document.getElementById('modal-name').value = button.getAttribute('data-name');
        document.getElementById('modal-invnumber').value = button.getAttribute('data-invnumber');
        document.getElementById('modal-hostelnumber').value = button.getAttribute('data-hostelnumber');
        document.getElementById('modal-roomnumber').value = button.getAttribute('data-roomnumber');
    });
});
