export enum HttpProtocol {
    HTTPS = 'https',
    HTTP = 'http'
}

export enum HttpMethod {
    GET = 'get',
    POST = 'post'
}

export interface HttpHeaders {
    [key: string]: string;
}

export type HttpData<D = { [key: string]: unknown }> = D;

export type HttpResponse<T = unknown> = {
    status: number;
    headers: HttpHeaders;
    cookies: HttpHeaders;
    message?: string;
    data?: T;
    success?: boolean;
};
