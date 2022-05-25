import { Identifiable } from '../../domain/interfaces/identifiable';
import { Ids } from '../../domain/definitions/key-types';

export abstract class BaseModel<T extends Identifiable> {
    public readonly fqid!: string;
    public readonly deleted!: boolean;
    public readonly data!: T &
        Identifiable & {
            meta_deleted: boolean;
            meta_position: number;
        };
}

export interface SearchResult<Model extends Identifiable> {
    collection: string;
    ids: Ids;
    fields: { [key in keyof Model]?: null };
}

export interface Repository<Model extends Identifiable> {
    /**
     * A unique collection to distinguish models
     */
    readonly COLLECTION: string;
    /**
     *
     */
    search(query: string): Promise<SearchResult<Model>>;
    /**
     *
     */
    buildColumn(): Promise<void>;
    /**
     *
     */
    buildIndex(): Promise<void>;
    /**
     *
     */
    buildTriggerFn(): Promise<void>;
}
