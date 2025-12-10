/**
 * app.js - HTMX-inspired library for SSE-based DOM updates
 * 
 * Attributes:
 *   data-sse-connect="<url>"  - Connect to SSE endpoint, updates element with received HTML
 *   data-sse-target="<selector>" - Optional target selector for SSE updates (default: self)
 *   data-search-input="<selector>" - Input triggers filtering on target element
 *   data-search-param="<name>" - URL parameter name for search (default: "q")
 *   data-filter-list - Mark a list as filterable
 *   data-filter-value="<text>" - Value to match against for filtering
 *   data-limit="<n>" - Maximum items to show
 */
(function() {
    'use strict';

    const App = {
        sseConnections: new Map(),
        
        init() {
            this.initSSE();
            this.initSearch();
            this.syncFromURL();
            window.addEventListener('popstate', () => this.syncFromURL());
        },

        // SSE Connection Management
        initSSE() {
            document.querySelectorAll('[data-sse-connect]').forEach(el => {
                this.connectSSE(el);
            });
        },

        connectSSE(el) {
            const url = el.dataset.sseConnect;
            if (!url) return;

            const targetSelector = el.dataset.sseTarget;
            const target = targetSelector ? document.querySelector(targetSelector) : el;
            
            if (this.sseConnections.has(el)) {
                this.sseConnections.get(el).close();
            }

            const eventSource = new EventSource(url);
            this.sseConnections.set(el, eventSource);

            eventSource.onmessage = (event) => {
                this.handleSSEMessage(target, event.data);
            };

            eventSource.addEventListener('update', (event) => {
                this.handleSSEMessage(target, event.data);
            });

            eventSource.addEventListener('append', (event) => {
                this.handleSSEAppend(target, event.data);
            });

            eventSource.addEventListener('remove', (event) => {
                const selector = event.data.trim();
                const toRemove = target.querySelector(selector);
                if (toRemove) toRemove.remove();
            });

            eventSource.onerror = () => {
                setTimeout(() => this.connectSSE(el), 5000);
            };
        },

        handleSSEMessage(target, html) {
            if (!target) return;
            target.innerHTML = html;
            this.applyFiltering(target);
        },

        handleSSEAppend(target, html) {
            if (!target) return;
            const template = document.createElement('template');
            template.innerHTML = html.trim();
            const newNodes = template.content.childNodes;
            newNodes.forEach(node => {
                if (node.nodeType === Node.ELEMENT_NODE) {
                    target.appendChild(node.cloneNode(true));
                }
            });
            this.applyFiltering(target);
        },

        // Search and Filtering
        initSearch() {
            document.querySelectorAll('[data-search-input]').forEach(input => {
                const targetSelector = input.dataset.searchInput;
                const paramName = input.dataset.searchParam || 'q';
                
                input.addEventListener('input', debounce(() => {
                    const value = input.value;
                    this.updateURL(paramName, value);
                    this.filterList(targetSelector, value);
                }, 150));
            });
        },

        syncFromURL() {
            const params = new URLSearchParams(window.location.search);
            
            document.querySelectorAll('[data-search-input]').forEach(input => {
                const paramName = input.dataset.searchParam || 'q';
                const value = params.get(paramName) || '';
                input.value = value;
                
                const targetSelector = input.dataset.searchInput;
                this.filterList(targetSelector, value);
            });
        },

        updateURL(param, value) {
            const url = new URL(window.location);
            if (value) {
                url.searchParams.set(param, value);
            } else {
                url.searchParams.delete(param);
            }
            window.history.replaceState({}, '', url);
        },

        filterList(targetSelector, query) {
            const target = document.querySelector(targetSelector);
            if (!target) return;
            
            this.applyFiltering(target, query);
        },

        applyFiltering(container, query) {
            if (!container) return;
            
            const filterList = container.closest('[data-filter-list]') || 
                               container.querySelector('[data-filter-list]') ||
                               (container.hasAttribute('data-filter-list') ? container : null);
            
            if (!filterList) return;
            
            const limit = parseInt(filterList.dataset.limit) || Infinity;
            const items = filterList.querySelectorAll('[data-filter-value]');
            
            // Get query from argument or from associated search input
            if (query === undefined) {
                const searchInput = document.querySelector(`[data-search-input="${getSelector(filterList)}"]`);
                query = searchInput ? searchInput.value : '';
            }
            
            const normalizedQuery = query.toLowerCase().trim();
            let visibleCount = 0;

            items.forEach(item => {
                const value = (item.dataset.filterValue || '').toLowerCase();
                const matches = fuzzyMatch(normalizedQuery, value);
                
                if (matches && visibleCount < limit) {
                    item.style.display = '';
                    item.classList.remove('filtered-out');
                    visibleCount++;
                } else {
                    item.style.display = 'none';
                    item.classList.add('filtered-out');
                }
            });

            // Update count display if exists
            const countEl = document.querySelector('[data-filter-count]');
            if (countEl) {
                const total = items.length;
                if (normalizedQuery || visibleCount < total) {
                    countEl.textContent = `${visibleCount} / ${total}`;
                } else {
                    countEl.textContent = `${total}`;
                }
            }
        }
    };

    // Fuzzy matching - matches if all characters appear in order
    function fuzzyMatch(pattern, text) {
        if (!pattern) return true;
        
        let patternIdx = 0;
        for (let i = 0; i < text.length && patternIdx < pattern.length; i++) {
            if (text[i] === pattern[patternIdx]) {
                patternIdx++;
            }
        }
        return patternIdx === pattern.length;
    }

    function debounce(fn, delay) {
        let timeout;
        return function(...args) {
            clearTimeout(timeout);
            timeout = setTimeout(() => fn.apply(this, args), delay);
        };
    }

    function getSelector(el) {
        if (el.id) return '#' + el.id;
        if (el.className) return '.' + el.className.split(' ')[0];
        return el.tagName.toLowerCase();
    }

    // Initialize on DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => App.init());
    } else {
        App.init();
    }

    // Expose for manual control
    window.App = App;
})();
