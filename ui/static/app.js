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
            const items = Array.from(filterList.querySelectorAll('[data-filter-value]'));
            
            // Get query from argument or from associated search input
            if (query === undefined) {
                // Find input that targets this filter list
                const listId = filterList.id;
                const searchInput = listId 
                    ? document.querySelector(`[data-search-input="#${listId}"]`)
                    : null;
                query = searchInput ? searchInput.value : '';
            }
            
            const normalizedQuery = query.toLowerCase().trim();
            
            // Score and sort items
            const scored = items.map(item => {
                const value = item.dataset.filterValue || '';
                const score = fuzzyScore(normalizedQuery, value.toLowerCase(), value);
                return { item, score };
            });
            
            // Sort by score (higher is better), then show top results
            scored.sort((a, b) => b.score - a.score);
            
            let visibleCount = 0;
            scored.forEach(({ item, score }) => {
                if (score > 0 && visibleCount < limit) {
                    item.style.display = '';
                    item.classList.remove('filtered-out');
                    visibleCount++;
                } else {
                    item.style.display = 'none';
                    item.classList.add('filtered-out');
                }
            });
            
            // Reorder DOM to match sorted order
            scored.forEach(({ item }) => {
                filterList.appendChild(item);
            });

            // Update count display if exists
            const countEl = document.querySelector('[data-filter-count]');
            if (countEl) {
                const total = items.length;
                const matchCount = scored.filter(s => s.score > 0).length;
                if (normalizedQuery) {
                    countEl.textContent = `${Math.min(visibleCount, matchCount)} / ${matchCount}`;
                } else {
                    countEl.textContent = `${visibleCount} / ${total}`;
                }
            }
        }
    };

    // Fuzzy scoring - returns score based on match quality (higher is better, 0 = no match)
    function fuzzyScore(pattern, textLower, textOriginal) {
        if (!pattern) return 1; // Empty pattern matches everything
        
        // Extract simple class name from fully qualified name
        const lastDot = textOriginal.lastIndexOf('.');
        const simpleName = lastDot >= 0 ? textOriginal.substring(lastDot + 1) : textOriginal;
        const simpleNameLower = simpleName.toLowerCase();
        
        // Exact match on simple name (highest priority)
        if (simpleNameLower === pattern) return 10000;
        
        // Simple name starts with pattern
        if (simpleNameLower.startsWith(pattern)) return 5000 + (pattern.length / simpleName.length) * 1000;
        
        // Simple name contains pattern as substring
        if (simpleNameLower.includes(pattern)) return 3000 + (pattern.length / simpleName.length) * 500;
        
        // Exact match on full name
        if (textLower === pattern) return 2000;
        
        // Full name starts with pattern
        if (textLower.startsWith(pattern)) return 1500;
        
        // Full name contains pattern as substring
        if (textLower.includes(pattern)) return 1000;
        
        // Fuzzy match - all characters appear in order
        let patternIdx = 0;
        let score = 0;
        let lastMatchIdx = -1;
        
        for (let i = 0; i < textLower.length && patternIdx < pattern.length; i++) {
            if (textLower[i] === pattern[patternIdx]) {
                // Bonus for consecutive matches
                if (lastMatchIdx === i - 1) score += 10;
                // Bonus for matching at word boundaries (after . or at start)
                if (i === 0 || textLower[i - 1] === '.') score += 20;
                score += 1;
                lastMatchIdx = i;
                patternIdx++;
            }
        }
        
        return patternIdx === pattern.length ? score : 0;
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
