function applyFilters() {
    let instituteFilter = document.getElementById("filterInstitute").value.toLowerCase();
    let rows = document.querySelectorAll(".resident-row");

    rows.forEach(row => {
        let institute = row.querySelector(".institute-row").textContent.toLowerCase();

        if (instituteFilter === "" || institute.includes(instituteFilter)) {
            row.style.display = "";
        } else {
            row.style.display = "none";
        }
    });
}

function generatePassword() {
    const length = 10;
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*";
    let password = "";
    for (let i = 0; i < length; i++) {
        password += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    document.getElementById("password").value = password;
}

document.getElementById("searchInput").addEventListener("input", function() {
    let filter = this.value.toLowerCase();
    let rows = document.querySelectorAll("#residentsTable tbody tr");

    rows.forEach(row => {
        let name = row.querySelector(".resident-name").textContent.toLowerCase();
        let room = row.querySelector(".room-number").textContent.toLowerCase();

        if (name.includes(filter) || room.includes(filter)) {
            row.style.display = "";
        } else {
            row.style.display = "none";
        }
    });
});