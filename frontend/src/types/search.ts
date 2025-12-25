export interface SearchResult {
  id: string;
  path: string;
  content: string;
  type: string;
  title?: string;
  _rankingScore?: number;
}

export interface SearchResponse {
  hits: SearchResult[];
  query: string;
  processingTimeMs: number;
  limit: number;
  offset: number;
  estimatedTotalHits: number;
}
