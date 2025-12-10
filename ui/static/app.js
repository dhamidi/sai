/**
 * app.js - Minimal search integration with Turbo Frames
 */
(function() {
    'use strict';

    function init() {
        const searchInput = document.getElementById('class-search');
        const classListFrame = document.getElementById('class-list');
        
        if (!searchInput || !classListFrame) return;

        let debounceTimer;
        
        searchInput.addEventListener('input', () => {
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => {
                const query = searchInput.value;
                const activeClass = getActiveClassName();
                
                let url = '/sidebar?q=' + encodeURIComponent(query);
                if (activeClass) {
                    url += '&active=' + encodeURIComponent(activeClass);
                }
                
                classListFrame.src = url;
            }, 150);
        });
    }

    function getActiveClassName() {
        const activeLink = document.querySelector('.class-item.active');
        if (!activeLink) return '';
        const href = activeLink.getAttribute('href');
        return href ? href.replace('/c/', '') : '';
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();
