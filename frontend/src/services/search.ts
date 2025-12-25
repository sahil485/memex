import { Search, GetMeilisearchHealth, OpenFile, IndexFile, IndexDirectory } from '../wailsjs/go/main/App';
import type { SearchResponse } from '../types/search';

export class SearchService {
  async search(query: string, limit: number = 20): Promise<SearchResponse> {
    const response = await Search(query, limit);

    if (response.error) {
      throw new Error(response.error);
    }

    return {
      hits: response.hits.map(hit => ({
        id: hit.id,
        path: hit.path,
        content: hit.content,
        type: hit.type,
        title: hit.title,
        _rankingScore: hit.rankingScore,
      })),
      query: response.query,
      processingTimeMs: response.processingTimeMs,
      limit: limit,
      offset: 0,
      estimatedTotalHits: response.estimatedTotalHits,
    };
  }

  async getHealth(): Promise<boolean> {
    try {
      return await GetMeilisearchHealth();
    } catch {
      return false;
    }
  }

  async openFile(path: string): Promise<void> {
    return OpenFile(path);
  }

  async indexFile(path: string): Promise<void> {
    return IndexFile(path);
  }

  async indexDirectory(path: string): Promise<void> {
    return IndexDirectory(path);
  }
}
