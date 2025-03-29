function applyFilters() {
    const searchTerm = document.getElementById("search").value.toLowerCase();
    const cards = document.querySelectorAll(".notice-card");
    let hasVisibleCards = false;

    cards.forEach(card => {
        const title = card.querySelector(".notice-title").textContent.toLowerCase();
        const annotation = card.querySelector(".notice-text").textContent.toLowerCase();

        if (searchTerm === "" || title.includes(searchTerm) || annotation.includes(searchTerm)) {
            card.style.display = "block";
            hasVisibleCards = true;
        } else {
            card.style.display = "none";
        }
    });

    const noResultsAlert = document.querySelector(".notice-cards-container .alert-info");
    if (noResultsAlert) {
        noResultsAlert.style.display = hasVisibleCards ? "none" : "block";
    }
}

function saveNotice() {
    // Логика сохранения объявления
    const modal = bootstrap.Modal.getInstance(document.getElementById('addNoticeModal'));
    modal.hide();
}