/**
 * app.js - Debounced search for Turbo Frames
 */
(function() {
    'use strict';

    function init() {
        const searchInput = document.getElementById('class-search');
        const searchForm = document.getElementById('search-form');

        if (!searchInput || !searchForm) return;

        let debounceTimer;

        searchInput.addEventListener('input', () => {
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => {
                searchForm.requestSubmit();
            }, 150);
        });
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
