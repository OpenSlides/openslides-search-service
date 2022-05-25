import { TestHttpService } from './test-http-service';

const SEARCH_SERVICE_URL = `http://localhost:9022/system/search`;

export class TestSearchServiceAdapter {
    public constructor(private readonly _http: TestHttpService) {}

    public search(query: string): Promise<any> {
        return this._http.post(SEARCH_SERVICE_URL, { data: { query } });
    }
}
