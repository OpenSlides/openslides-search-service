import { TestHttpService } from './test-http-service';
import { TEST_ADMIN_ID } from './config';
import { TestPostgreAdapter } from './test-postgre-adapter';

export enum EventType {
    CREATE = 'create',
    UPDATE = 'update',
    DELETE = 'delete'
}

export type DatastoreEventType = EventType;

abstract class Event {
    fqid!: string;
    abstract type: EventType;
}

class CreateEvent<T> extends Event {
    public readonly type = EventType.CREATE;
    fields!: { [key in keyof T]?: unknown };
}

class UpdateEvent<T> extends Event {
    public readonly type = EventType.UPDATE;
    fields!: { [key in keyof T]?: unknown };
}

class DeleteEvent extends Event {
    public readonly type = EventType.DELETE;
}

type DatastoreEvent<T = any> = CreateEvent<T> | UpdateEvent<T> | DeleteEvent;

class DatastoreSchema {
    public readonly user_id!: number;
    public readonly locked_fields!: object;
    public readonly information!: object;
    public readonly events!: DatastoreEvent[];
}

export class TestDatastoreWriterAdapter {
    public constructor(private readonly _postgre: TestPostgreAdapter, private readonly _http: TestHttpService) {}

    public async isReady(): Promise<boolean> {
        try {
            await this._postgre.openConnection();
            return true;
        } catch (e) {
            return false;
        }
    }

    public async prune(): Promise<void> {
        await this._postgre.prune();
    }

    public async write<T>(events: DatastoreEvent[]): Promise<T> {
        const url = this.getWriterUrl();
        const answer = await this._http.post<T, DatastoreSchema>(url, {
            data: {
                user_id: TEST_ADMIN_ID,
                information: {},
                locked_fields: {},
                events
            }
        });
        return answer.data!;
    }

    public async closeConnection(): Promise<void> {
        await this._postgre.closeConnection();
    }

    private getWriterUrl(): string {
        const writerHost = process.env.DATASTORE_WRITER_HOST;
        const writerPort = process.env.DATASTORE_WRITER_PORT;
        if (!writerHost || !writerPort) {
            throw new Error('No datastore writer is defined.');
        }
        return `http://${writerHost}:${parseInt(writerPort, 10)}/internal/datastore/writer/write`;
    }
}
