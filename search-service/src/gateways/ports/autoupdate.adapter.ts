import { Inject } from 'final-di';
import { HttpClient, HttpService } from '../http';
import { AutoupdatePort, AutoupdateRequestConfig } from './autoupdate.port';

const AUTOUPDATE_URL = process.env.AUTOUPDATE_URL || `http://localhost:9012/system/autoupdate`;

export class AutoupdateAdapter implements AutoupdatePort {
    @Inject(HttpService)
    private readonly _http!: HttpClient;

    public async request(data: AutoupdateRequestConfig<any>[]): Promise<{ [index: string]: any }> {
        return await this._http.post({ url: `${AUTOUPDATE_URL}?single=1`, data });
    }
}
