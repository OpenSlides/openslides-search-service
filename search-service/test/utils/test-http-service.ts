import { HttpClient, HttpService, HttpHeaders, HttpResponse } from '../../src/gateways/http';

export interface RequestOptions<D = any> {
    headers?: HttpHeaders;
    data?: D;
}

export class TestHttpService {
    private http: HttpClient = new HttpService();

    public async post<T = any, D = any>(url: string, options: RequestOptions<D> = {}): Promise<HttpResponse<T>> {
        return this.http.post<T, D>({ url, ...options });
    }
}
