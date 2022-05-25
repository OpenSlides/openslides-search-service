import { SearchResult } from '../repositories/repository';
export interface SearchClientResponse extends SearchResult<any> {}

export interface SearchClient {
    search(query: string): Promise<SearchClientResponse[]>;
}
