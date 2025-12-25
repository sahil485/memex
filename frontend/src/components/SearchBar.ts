import { SearchService } from '../services/search';
import type { SearchResult } from '../types/search';
import { Quit } from '../wailsjs/runtime/runtime';

export class SearchBar {
  private container: HTMLElement;
  private searchInput: HTMLInputElement;
  private resultsContainer: HTMLElement;
  private searchService: SearchService;
  private selectedIndex: number = 0;
  private results: SearchResult[] = [];
  private searchTimeout?: number;

  constructor() {
    this.searchService = new SearchService();
    this.container = this.createSearchBar();
    this.searchInput = this.container.querySelector('#search-input') as HTMLInputElement;
    this.resultsContainer = this.container.querySelector('#search-results') as HTMLElement;

    this.setupEventListeners();
    document.body.appendChild(this.container);

    // Focus input immediately
    setTimeout(() => this.searchInput.focus(), 100);
  }

  private createSearchBar(): HTMLElement {
    const container = document.createElement('div');
    container.style.cssText = `
      position: fixed;
      top: 60px;
      left: 50%;
      transform: translateX(-50%);
      width: 600px;
      max-width: 90vw;
      background: transparent;
      --wails-draggable: drag;
    `;

    container.innerHTML = `
      <div class="search-header" style="
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 14px 18px;
        background: rgba(255, 255, 255, 0.98);
        border-radius: 12px;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
        backdrop-filter: blur(40px);
        border: 1px solid rgba(0, 0, 0, 0.08);
        --wails-draggable: drag;
      ">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="none" style="flex-shrink: 0; opacity: 0.4;">
          <circle cx="7" cy="7" r="5.5" stroke="currentColor" stroke-width="1.5"/>
          <path d="M11 11L14.5 14.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        </svg>
        <input
          type="text"
          id="search-input"
          placeholder="Search files..."
          autocomplete="off"
          spellcheck="false"
          style="
            flex: 1;
            border: none;
            outline: none;
            font-size: 15px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            color: #111827;
            background: transparent;
            -webkit-app-region: no-drag;
          "
        />
        <button id="close-btn" style="
          flex-shrink: 0;
          width: 20px;
          height: 20px;
          border: none;
          background: transparent;
          cursor: pointer;
          display: flex;
          align-items: center;
          justify-content: center;
          border-radius: 4px;
          opacity: 0.5;
          transition: all 0.15s;
          -webkit-app-region: no-drag;
        " onmouseover="this.style.opacity='1'; this.style.background='rgba(0,0,0,0.05)'" onmouseout="this.style.opacity='0.5'; this.style.background='transparent'">
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
            <path d="M1 1L11 11M1 11L11 1" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
      </div>
      <div id="search-results-wrapper" style="
        display: none;
        margin-top: 8px;
        background: rgba(255, 255, 255, 0.98);
        border-radius: 12px;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
        backdrop-filter: blur(40px);
        border: 1px solid rgba(0, 0, 0, 0.08);
        max-height: 400px;
        overflow-y: auto;
      ">
        <div id="search-results"></div>
      </div>
    `;

    return container;
  }

  private setupEventListeners(): void {
    this.searchInput.addEventListener('input', () => {
      clearTimeout(this.searchTimeout);
      this.searchTimeout = window.setTimeout(() => this.performSearch(), 300);
    });

    this.searchInput.addEventListener('keydown', (e) => this.handleKeyDown(e));

    // Close button
    const closeBtn = this.container.querySelector('#close-btn');
    closeBtn?.addEventListener('click', () => {
      Quit();
    });
  }

  private handleKeyDown(e: KeyboardEvent): void {
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        this.selectNext();
        break;
      case 'ArrowUp':
        e.preventDefault();
        this.selectPrevious();
        break;
      case 'Enter':
        e.preventDefault();
        this.openSelected();
        break;
      case 'Escape':
        e.preventDefault();
        this.clearSearch();
        break;
    }
  }

  private async performSearch(): Promise<void> {
    const query = this.searchInput.value.trim();
    const wrapper = this.container.querySelector('#search-results-wrapper') as HTMLElement;

    if (!query) {
      wrapper.style.display = 'none';
      this.results = [];
      return;
    }

    try {
      const response = await this.searchService.search(query);
      this.results = response.hits;
      this.selectedIndex = 0;
      this.renderResults();
    } catch (error) {
      console.error('Search error:', error);
      this.resultsContainer.innerHTML = `
        <div style="padding: 20px; text-align: center; color: #ef4444;">
          Search error. Check if MeiliSearch is running.
        </div>
      `;
      wrapper.style.display = 'block';
    }
  }

  private renderResults(): void {
    const wrapper = this.container.querySelector('#search-results-wrapper') as HTMLElement;

    if (this.results.length === 0) {
      wrapper.style.display = 'none';
      return;
    }

    wrapper.style.display = 'block';
    this.resultsContainer.innerHTML = this.results
      .map((result, index) => this.renderResultItem(result, index))
      .join('');

    this.resultsContainer.querySelectorAll('.search-result-item').forEach((item, index) => {
      item.addEventListener('click', () => {
        this.selectedIndex = index;
        this.openSelected();
      });
    });
  }

  private renderResultItem(result: SearchResult, index: number): string {
    const isSelected = index === this.selectedIndex;
    const fileType = this.getFileType(result.path);
    const fileName = result.path.split('/').pop() || result.path;
    const filePath = result.path;

    return `
      <div class="search-result-item" data-index="${index}" style="
        padding: 10px 16px;
        cursor: pointer;
        border-bottom: 1px solid rgba(229, 231, 235, 0.5);
        background: ${isSelected ? 'rgba(240, 249, 255, 0.8)' : 'transparent'};
        transition: background 0.1s ease;
      ">
        <div style="display: flex; align-items: center; gap: 10px;">
          <div style="
            flex-shrink: 0;
            width: 28px;
            height: 28px;
            background: ${this.getFileTypeColor(fileType)};
            border-radius: 5px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 10px;
            font-weight: 600;
            color: #ffffff;
          ">
            ${fileType}
          </div>
          <div style="flex: 1; min-width: 0;">
            <div style="
              font-size: 13px;
              font-weight: 500;
              color: #111827;
              white-space: nowrap;
              overflow: hidden;
              text-overflow: ellipsis;
            ">
              ${fileName}
            </div>
            <div style="
              font-size: 11px;
              color: #6b7280;
              white-space: nowrap;
              overflow: hidden;
              text-overflow: ellipsis;
            ">
              ${filePath}
            </div>
          </div>
        </div>
      </div>
    `;
  }

  private getFileType(path: string): string {
    const ext = path.split('.').pop()?.toLowerCase() || '';
    const typeMap: Record<string, string> = {
      'js': 'JS', 'ts': 'TS', 'jsx': 'JSX', 'tsx': 'TSX',
      'py': 'PY', 'go': 'GO', 'java': 'JAVA',
      'cpp': 'C++', 'c': 'C', 'rs': 'RS',
      'md': 'MD', 'json': 'JSON', 'yaml': 'YML',
      'html': 'HTML', 'css': 'CSS', 'scss': 'SCSS',
    };
    return typeMap[ext] || ext.toUpperCase().slice(0, 3) || 'FILE';
  }

  private getFileTypeColor(type: string): string {
    const colorMap: Record<string, string> = {
      'JS': '#f7df1e', 'TS': '#3178c6', 'JSX': '#61dafb', 'TSX': '#3178c6',
      'PY': '#3776ab', 'GO': '#00add8', 'JAVA': '#007396',
      'C++': '#00599c', 'C': '#555555', 'RS': '#ce422b',
      'MD': '#083fa1', 'JSON': '#5a5a5a', 'YML': '#cb171e',
      'HTML': '#e34c26', 'CSS': '#1572b6', 'SCSS': '#cc6699',
    };
    return colorMap[type] || '#6b7280';
  }

  private selectNext(): void {
    if (this.results.length === 0) return;
    this.selectedIndex = (this.selectedIndex + 1) % this.results.length;
    this.renderResults();
    this.scrollToSelected();
  }

  private selectPrevious(): void {
    if (this.results.length === 0) return;
    this.selectedIndex = this.selectedIndex === 0
      ? this.results.length - 1
      : this.selectedIndex - 1;
    this.renderResults();
    this.scrollToSelected();
  }

  private scrollToSelected(): void {
    const selectedElement = this.resultsContainer.querySelector(
      `[data-index="${this.selectedIndex}"]`
    );
    selectedElement?.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
  }

  private async openSelected(): Promise<void> {
    if (this.results.length === 0) return;
    const selected = this.results[this.selectedIndex];

    try {
      await this.searchService.openFile(selected.path);
      console.log('Opened file:', selected.path);
    } catch (error) {
      console.error('Error opening file:', error);
    }
  }

  private clearSearch(): void {
    this.searchInput.value = '';
    this.results = [];
    const wrapper = this.container.querySelector('#search-results-wrapper') as HTMLElement;
    wrapper.style.display = 'none';
  }
}
