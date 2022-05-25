import { Inject } from 'final-di';
import { BaseModel, Repository, SearchResult } from './repository';
import { PostgreAdapterService, PostgreService } from '../postgre';
import { Identifiable } from '../../domain/interfaces/identifiable';

export abstract class SearchRepository<Model extends Identifiable = any> implements Repository<Model> {
    public abstract readonly COLLECTION: string;

    @Inject(PostgreAdapterService)
    private readonly _postgre!: PostgreService;

    public async search(query: string): Promise<SearchResult<Model>> {
        const result = await this._postgre.select<BaseModel<Model>>(this.getCollection(), query);
        return {
            collection: this.COLLECTION,
            ids: result.rows.map(entry => entry.data.id),
            fields: this.getSearchableFields().mapToObject(entry => ({ [entry]: null }))
        };
    }

    public buildColumn(): Promise<void> {
        return this._postgre.createColumn(this.getCollection());
    }

    public buildIndex(): Promise<void> {
        return this._postgre.createIndex(this.getCollection(), this.getSearchableFields() as string[]);
    }

    public buildTriggerFn(): Promise<void> {
        return this._postgre.createTriggerFn(this.getCollection(), this.getSearchableFields() as string[]);
    }

    public getCollection(): string {
        return this.COLLECTION;
    }

    public abstract getSearchableFields(): (keyof Model)[];
}
