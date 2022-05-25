import { SearchResult } from '../repositories/repository';
import { Identifiable } from '../../domain/interfaces/identifiable';

export interface AutoupdateRequestConfig<T extends Identifiable = any> extends SearchResult<T> {}

export interface AutoupdatePort {
    request(data: AutoupdateRequestConfig[]): Promise<{ [index: string]: any }>;
}
