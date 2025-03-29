function applyFilters() {
    const searchTerm = document.getElementById("searchNews").value.toLowerCase();
    const newsCards = document.querySelectorAll(".news-card");
    
    let hasVisibleCards = false;

    newsCards.forEach(card => {
        const title = card.querySelector(".news-title").textContent.toLowerCase();
        const annotation = card.querySelector(".news-annotation").textContent.toLowerCase();
        
        if (searchTerm === "" || title.includes(searchTerm) || annotation.includes(searchTerm)) {
            card.style.display = "block";
            hasVisibleCards = true;
        } else {
            card.style.display = "none";
        }
    });

    const noResultsAlert = document.querySelector(".news-cards-container .alert-info");
    if (noResultsAlert) {
        noResultsAlert.style.display = hasVisibleCards ? "none" : "block";
    }
}

function saveNews() {
    // Логика сохранения новости
    const modal = bootstrap.Modal.getInstance(document.getElementById('addNewsModal'));
    modal.hide();
}