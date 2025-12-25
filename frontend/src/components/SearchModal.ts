import { SearchService } from '../services/search';
import type { SearchResult } from '../types/search';

export class SearchModal {
  private modal: HTMLElement;
  private overlay: HTMLElement;
  private searchInput: HTMLInputElement;
  private resultsContainer: HTMLElement;
  private searchService: SearchService;
  private selectedIndex: number = 0;
  private results: SearchResult[] = [];
  private searchTimeout?: number;

  constructor() {
    this.searchService = new SearchService();
    this.modal = this.createModal();
    this.overlay = this.createOverlay();
    this.searchInput = this.modal.querySelector('#search-input') as HTMLInputElement;
    this.resultsContainer = this.modal.querySelector('#search-results') as HTMLElement;

    this.setupEventListeners();
    this.appendToDOM();
  }

  private createOverlay(): HTMLElement {
    const overlay = document.createElement('div');
    overlay.id = 'search-overlay';
    overlay.className = 'search-overlay';
    overlay.style.cssText = `
      display: none;
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: transparent;
      z-index: 999;
    `;
    return overlay;
  }

  private createModal(): HTMLElement {
    const modal = document.createElement('div');
    modal.id = 'search-modal';
    modal.className = 'search-modal';
    modal.style.cssText = `
      display: none;
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: transparent;
    `;

    modal.innerHTML = `
      <div class="search-header" style="
        padding: 12px 16px;
        background: rgba(255, 255, 255, 0.95);
        border-radius: 10px;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
        backdrop-filter: blur(10px);
      ">
        <input
          type="text"
          id="search-input"
          placeholder="Search..."
          autocomplete="off"
          spellcheck="false"
          style="
            width: 100%;
            border: none;
            outline: none;
            font-size: 16px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            color: #111827;
            background: transparent;
          "
        />
      </div>
      <div class="search-body" style="
        display: none;
        max-height: 400px;
        overflow-y: auto;
        background: rgba(255, 255, 255, 0.95);
        margin-top: 4px;
        border-radius: 10px;
        box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
        backdrop-filter: blur(10px);
      ">
        <div id="search-results"></div>
      </div>
    `;

    return modal;
  }

  private setupEventListeners(): void {
    // Search input
    this.searchInput.addEventListener('input', () => {
      clearTimeout(this.searchTimeout);
      this.searchTimeout = window.setTimeout(() => this.performSearch(), 300);
    });

    // Keyboard navigation
    this.searchInput.addEventListener('keydown', (e) => this.handleKeyDown(e));

    // Close on overlay click
    this.overlay.addEventListener('click', () => this.close());

    // Global keyboard shortcut (Cmd+K or Ctrl+K)
    document.addEventListener('keydown', (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        this.toggle();
      }
    });
  }

  private handleKeyDown(e: KeyboardEvent): void {
    switch (e.key) {
      case 'Escape':
        e.preventDefault();
        this.close();
        break;
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
    }
  }

  private async performSearch(): Promise<void> {
    const query = this.searchInput.value.trim();

    if (!query) {
      this.resultsContainer.innerHTML = '';
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
        <div style="padding: 40px 20px; text-align: center; color: #ef4444;">
          Error performing search. Make sure MeiliSearch is running.
        </div>
      `;
    }
  }

  private renderResults(): void {
    const searchBody = this.modal.querySelector('.search-body') as HTMLElement;

    if (this.results.length === 0) {
      searchBody.style.display = 'none';
      return;
    }

    searchBody.style.display = 'block';
    this.resultsContainer.innerHTML = this.results
      .map((result, index) => this.renderResultItem(result, index))
      .join('');

    // Add click handlers
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
        padding: 12px 20px;
        cursor: pointer;
        border-bottom: 1px solid #f3f4f6;
        background: ${isSelected ? '#f0f9ff' : '#ffffff'};
        transition: background 0.15s ease;
      ">
        <div style="display: flex; align-items: center; gap: 12px;">
          <div class="file-icon" style="
            flex-shrink: 0;
            width: 32px;
            height: 32px;
            background: ${this.getFileTypeColor(fileType)};
            border-radius: 6px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 11px;
            font-weight: 600;
            color: #ffffff;
            text-transform: uppercase;
          ">
            ${fileType}
          </div>
          <div style="flex: 1; min-width: 0;">
            <div style="
              font-size: 14px;
              font-weight: 500;
              color: #111827;
              margin-bottom: 4px;
              white-space: nowrap;
              overflow: hidden;
              text-overflow: ellipsis;
            ">
              ${fileName}
            </div>
            <div style="
              font-size: 12px;
              color: #6b7280;
              white-space: nowrap;
              overflow: hidden;
              text-overflow: ellipsis;
            ">
              ${filePath}
            </div>
          </div>
          ${result._rankingScore ? `
            <div style="
              flex-shrink: 0;
              font-size: 11px;
              color: #9ca3af;
              font-family: monospace;
            ">
              ${(result._rankingScore * 100).toFixed(0)}%
            </div>
          ` : ''}
        </div>
      </div>
    `;
  }

  private getFileType(path: string): string {
    const ext = path.split('.').pop()?.toLowerCase() || '';
    const typeMap: Record<string, string> = {
      'js': 'JS',
      'ts': 'TS',
      'jsx': 'JSX',
      'tsx': 'TSX',
      'py': 'PY',
      'go': 'GO',
      'java': 'JAVA',
      'cpp': 'C++',
      'c': 'C',
      'rs': 'RS',
      'md': 'MD',
      'json': 'JSON',
      'yaml': 'YAML',
      'yml': 'YAML',
      'html': 'HTML',
      'css': 'CSS',
      'scss': 'SCSS',
    };
    return typeMap[ext] || ext.toUpperCase().slice(0, 4) || 'FILE';
  }

  private getFileTypeColor(type: string): string {
    const colorMap: Record<string, string> = {
      'JS': '#f7df1e',
      'TS': '#3178c6',
      'JSX': '#61dafb',
      'TSX': '#3178c6',
      'PY': '#3776ab',
      'GO': '#00add8',
      'JAVA': '#007396',
      'C++': '#00599c',
      'C': '#555555',
      'RS': '#ce422b',
      'MD': '#083fa1',
      'JSON': '#5a5a5a',
      'YAML': '#cb171e',
      'HTML': '#e34c26',
      'CSS': '#1572b6',
      'SCSS': '#cc6699',
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
      // Use Wails to open the file natively
      await this.searchService.openFile(selected.path);
      console.log('Opened file:', selected.path);
      this.close();
    } catch (error) {
      console.error('Error opening file:', error);
      // Show error message in UI
      alert(`Failed to open file: ${selected.path}`);
    }
  }

  private appendToDOM(): void {
    document.body.appendChild(this.overlay);
    document.body.appendChild(this.modal);
  }

  open(): void {
    this.overlay.style.display = 'block';
    this.modal.style.display = 'block';
    this.searchInput.value = '';
    this.resultsContainer.innerHTML = '';
    this.results = [];
    this.selectedIndex = 0;

    // Focus input after a brief delay to ensure modal is visible
    setTimeout(() => this.searchInput.focus(), 50);
  }

  close(): void {
    this.overlay.style.display = 'none';
    this.modal.style.display = 'none';
  }

  toggle(): void {
    if (this.modal.style.display === 'none') {
      this.open();
    } else {
      this.close();
    }
  }
}
