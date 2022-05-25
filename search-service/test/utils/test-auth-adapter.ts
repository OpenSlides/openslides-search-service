import { TEST_ADMIN_USERNAME, TEST_ADMIN_PASSWORD } from './config';
import { TestHttpService } from './test-http-service';

const AUTH_URL = process.env.AUTH_URL || 'http://localhost:9004/system/auth';

const TestUser = {
    accessToken: '',
    refreshId: ''
};

export class TestAuthAdapter {
    public constructor(private readonly _http: TestHttpService) {}

    public async login(username: string = TEST_ADMIN_USERNAME, password: string = TEST_ADMIN_PASSWORD): Promise<void> {
        const response = await this._http.post<void>(`${AUTH_URL}/login`, { data: { username, password } });
        TestUser.accessToken = response.headers[`authentication`];
        TestUser.refreshId = response.cookies[`refreshId`];
    }
}
