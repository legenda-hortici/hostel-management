document.addEventListener("DOMContentLoaded", function () {
    // Обработчик для панели управления баннерами
    const toggleBannerSettingsBtn = document.getElementById('toggleBannerSettings');
    const bannerSettings = document.getElementById('bannerSettings');
    
    if (toggleBannerSettingsBtn && bannerSettings) {
        toggleBannerSettingsBtn.addEventListener('click', function() {
            const isVisible = bannerSettings.style.display === 'block';
            bannerSettings.style.display = isVisible ? 'none' : 'block';
        });
    }

    // Существующий код для каруселей
    const newsCarousel = document.getElementById("newsCarousel");
    const noticesCarousel = document.getElementById("noticesCarousel");
    const newsPrevBtn = document.getElementById("newsPrev");
    const newsNextBtn = document.getElementById("newsNext");
    const noticesPrevBtn = document.getElementById("noticesPrev");
    const noticesNextBtn = document.getElementById("noticesNext");

    const scrollAmount = 320;

    newsPrevBtn.addEventListener("click", function () {
        newsCarousel.scrollBy({ left: -scrollAmount, behavior: "smooth" });
    });

    newsNextBtn.addEventListener("click", function () {
        newsCarousel.scrollBy({ left: scrollAmount, behavior: "smooth" });
    });

    noticesPrevBtn.addEventListener("click", function () {
        noticesCarousel.scrollBy({ left: -scrollAmount, behavior: "smooth" });
    });

    noticesNextBtn.addEventListener("click", function () {
        noticesCarousel.scrollBy({ left: scrollAmount, behavior: "smooth" });
    });
});