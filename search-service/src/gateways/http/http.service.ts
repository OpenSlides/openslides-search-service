import axios, { AxiosError, AxiosResponse } from 'axios';
import { Logger } from '../../infrastructure/utils/logger';
import { HttpResponse, HttpMethod, HttpHeaders, HttpData } from './definitions';
import { HttpClient, HttpRequestConfig } from './http-client';

const DEFAULT_HEADERS = { 'Content-Type': 'application/json', Accept: 'application/json' };

export class HttpService implements HttpClient {
    public post<R, D>(config: HttpRequestConfig<D>): Promise<HttpResponse<R>> {
        return this.send<R, D>({ method: HttpMethod.POST, ...config });
    }

    public async send<R, D>(config: { method: HttpMethod } & HttpRequestConfig<D>): Promise<HttpResponse<R>> {
        const { url, data, headers = DEFAULT_HEADERS, method } = config;
        Logger.debug(`Sending a request: ${method} ${url} ${JSON.stringify(data)} ${JSON.stringify(headers)}`);
        Logger.debug(`${url} -H '${JSON.stringify(headers)}' -d '${JSON.stringify(data)}'`);
        try {
            const response = await axios({ url, method, data, headers, responseType: 'json' });
            return this.createHttpResponse<R>(response);
        } catch (e) {
            const message = (e as AxiosError).message;
            this.handleError(message, url, method, data, headers);
            return this.createHttpResponse<R>((e as AxiosError).response as AxiosResponse);
        }
    }

    private handleError(
        error: string,
        url: string,
        method: HttpMethod,
        data?: HttpData<any>,
        headers?: HttpHeaders
    ): void {
        Logger.error('HTTP-error occurred: ', error);
        Logger.error(`Error is occurred while sending the following information: ${method} ${url}`);
        Logger.error(
            `Request contains the following data ${JSON.stringify(data)} and headers ${JSON.stringify(headers)}`
        );
    }

    private createHttpResponse<T>(response: AxiosResponse<T>): HttpResponse<T> {
        const result = {
            status: response.status,
            headers: response.headers as HttpHeaders,
            cookies: this.getCookiesByHeaders(response.headers)
        };
        return {
            ...result,
            ...response.data
        };
    }

    private getCookiesByHeaders(headers: HttpHeaders): HttpHeaders {
        const parseCookie = (rawCookie: string): [string, string] => {
            const indexOfEqual = rawCookie.indexOf('=');
            const parts = [rawCookie.slice(0, indexOfEqual), rawCookie.slice(indexOfEqual + 1)];
            const pathIndex = parts[1].search(/Path=/i);
            return [parts[0], parts[1].slice(0, pathIndex - 2)];
        };
        const rawCookies = headers['set-cookie'] as any;
        const cookies: { [cookieName: string]: string } = {};
        if (rawCookies && rawCookies.length) {
            for (const rawCookie of rawCookies) {
                const [cookieKey, cookieValue]: [string, string] = parseCookie(rawCookie);
                cookies[cookieKey] = cookieValue;
            }
        }
        return cookies;
    }
}
