function applyFilters() {
    const filterSelect = document.getElementById("filterInstitute");
    const rows = document.querySelectorAll("#residentsTable tbody tr");
    const selectedFurniture = filterSelect.value.toLowerCase();

    rows.forEach(row => {
        const furnitureType = row.querySelector(".inv-name").textContent.toLowerCase();

        if (selectedFurniture === "" || furnitureType === selectedFurniture) {
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