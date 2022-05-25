import { TestDatastoreWriterAdapter } from './test-datastore-writer-adapter';
import { TestHttpService } from './test-http-service';
import { TestPostgreAdapter } from './test-postgre-adapter';
import { TestSearchServiceAdapter } from './test-search-service-adapter';
import { TestUserRepositoryAdapter } from '../suites/users/test-user-repository-adapter';
import { TestAuthAdapter } from './test-auth-adapter';

export class TestContainer {
    public readonly postgre: TestPostgreAdapter;
    public readonly http: TestHttpService;
    public readonly datastore: TestDatastoreWriterAdapter;
    public readonly search: TestSearchServiceAdapter;
    public readonly auth: TestAuthAdapter;

    private readonly _userRepo: TestUserRepositoryAdapter;

    public constructor() {
        this.postgre = new TestPostgreAdapter();
        this.http = new TestHttpService();
        this.datastore = new TestDatastoreWriterAdapter(this.postgre, this.http);
        this.search = new TestSearchServiceAdapter(this.http);
        this.auth = new TestAuthAdapter(this.http);
        this._userRepo = new TestUserRepositoryAdapter(this.datastore);
    }

    /**
     * Ensures that a connection to the database is established
     */
    public async open(): Promise<void> {
        await this.postgre.openConnection();
    }

    /**
     * Closes the connection to the database
     */
    public async close(): Promise<void> {
        await this.postgre.closeConnection();
    }

    public async init(): Promise<void> {
        await this.postgre.prune();
        await this._userRepo.createSuperadmin();
    }
}
