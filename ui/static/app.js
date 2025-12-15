/**
 * app.js - Debounced search for Turbo Frames
 */
(function() {
    'use strict';

    let debounceTimer;

    function init() {
        const searchInput = document.getElementById('class-search');
        const searchForm = document.getElementById('search-form');

        if (!searchInput || !searchForm) return;

        // Avoid adding duplicate listeners
        if (searchInput.dataset.initialized) return;
        searchInput.dataset.initialized = 'true';

        searchInput.addEventListener('input', () => {
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => {
                searchForm.requestSubmit();
            }, 150);
        });
    }

    // Run init immediately if DOM is ready (covers deferred script case)
    if (document.readyState !== 'loading') {
        init();
    } else {
        document.addEventListener('DOMContentLoaded', init);
    }

    // Also run after Turbo navigations
    document.addEventListener('turbo:load', init);
})();
