import { HttpData, HttpHeaders, HttpMethod, HttpResponse } from './definitions';

export interface HttpRequestConfig<D = { [index: string]: unknown }> {
    url: string;
    data?: HttpData<D>;
    headers?: HttpHeaders;
}

export interface HttpClient {
    post<R, D>(requestConfig: HttpRequestConfig<D>): Promise<HttpResponse<R>>;
    send<R, D>(requestConfig: { method: HttpMethod } & HttpRequestConfig<D>): Promise<HttpResponse<R>>;
}
