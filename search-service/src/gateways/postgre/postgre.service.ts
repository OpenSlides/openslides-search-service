import { Client, QueryResult } from 'pg';

export interface PostgreService {
    getFqids(): Promise<{ fqid: string }[]>;
    select<R = any>(columnName: string, query: string): Promise<QueryResult<R>>;
    createColumn(columnName: string): Promise<void>;
    createIndex(name: string, indexedFields: string[]): Promise<void>;
    createTriggerFn(name: string, indexedFields: string[]): Promise<void>;
    getPgClient(): Promise<Client | null>;
}
