import { EventType, TestDatastoreWriterAdapter } from './test-datastore-writer-adapter';
import { Identifiable } from '../../src/domain/interfaces/identifiable';
import { Fqid } from '../../src/domain/definitions/key-types';
import { Mutable } from '../../src/domain/interfaces/mutable';

let uniqueId = 1;

export abstract class RepositoryAdapter<Model extends Mutable<Identifiable>> {
    public constructor(private readonly _datastore: TestDatastoreWriterAdapter) {}

    public async create(model: Partial<Model>): Promise<number> {
        const nextId = uniqueId++;
        const fqid = this.getFqid(nextId);
        model.id = model.id || nextId;
        await this._datastore.write([{ type: EventType.CREATE, fqid, fields: model }]);
        return nextId;
    }

    public async update(id: number, model: Partial<Model>): Promise<void> {
        const fqid = this.getFqid(id);
        await this._datastore.write([{ type: EventType.UPDATE, fqid, fields: model }]);
    }

    public async delete(id: number): Promise<void> {
        const fqid = this.getFqid(id);
        await this._datastore.write([{ type: EventType.DELETE, fqid }]);
    }

    protected abstract getFqid(id: number): Fqid;
}
