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
                
                fetch(url)
                    .then(r => r.text())
                    .then(html => {
                        const parser = new DOMParser();
                        const doc = parser.parseFromString(html, 'text/html');
                        const newFrame = doc.getElementById('class-list');
                        if (newFrame) {
                            classListFrame.innerHTML = newFrame.innerHTML;
                        }
                    });
            }, 150);
        });

        document.addEventListener('click', (e) => {
            const link = e.target.closest('.class-item');
            if (!link) return;
            
            document.querySelectorAll('.class-item.active').forEach(el => {
                el.classList.remove('active');
            });
            link.classList.add('active');
            
            const href = link.getAttribute('href');
            if (href) {
                history.pushState(null, '', href);
            }
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
