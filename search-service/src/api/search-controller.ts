import { Inject } from 'final-di';
import { Body, OnGet, RestController, Cookie, Header, RestMiddleware, RoutingError, OnPost } from 'rest-app';
import { Logger } from '../infrastructure/utils/logger';
import { SearchClient } from '../gateways/ports/search-client';
import { SearchService } from '../gateways/ports/search.service';
import { RestServerResponse } from '../infrastructure/utils/definitions';
import { createResponse } from '../infrastructure/utils/functions';
import { AutoupdateAdapter } from '../gateways/ports/autoupdate.adapter';
import { AutoupdatePort } from '../gateways/ports/autoupdate.port';
import { Collection } from '../domain/definitions/key-types';
import { SEARCH_SERVICE_HOST, SEARCH_SERVICE_PORT } from '../infrastructure/utils/config';
import { Request, Response, NextFunction } from 'express';

interface SearchServiceRequest {
    query: string;
    collections?: Collection[];
}

@RestController({ prefix: `system/search` })
export class SearchController {
    @Inject(SearchService)
    private readonly _searchClient!: SearchClient;

    @Inject(AutoupdateAdapter)
    private readonly _autoupdate!: AutoupdatePort;

    @OnPost()
    public async index(
        @Body() data: SearchServiceRequest,
        @Header('authentication') accessToken: string,
        @Cookie('refreshId') refreshCookie: string
    ): Promise<RestServerResponse> {
        if (!isSearchServiceRequest(data)) {
            throw new RoutingError(`Data must contain "query" and can have "collections".`, { statusCode: 400 });
        }
        const result = await this._searchClient.search(data.query);
        Logger.debug(`Found result:`, JSON.stringify(result));
        const autoupdate = await this._autoupdate.request(result);
        Logger.debug(`Sending result:`, autoupdate);
        return createResponse({ message: 'Search could not found anything', results: [result] });
    }

    @OnGet()
    public health(): RestServerResponse {
        return createResponse({
            message: `search-service is available under ${SEARCH_SERVICE_HOST}:${SEARCH_SERVICE_PORT}`
        });
    }
}

function isSearchServiceRequest(request: any): request is SearchServiceRequest {
    return typeof request?.query === `string` && request.query.length >= 1;
}
