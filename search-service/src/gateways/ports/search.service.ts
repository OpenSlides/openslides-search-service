import { Inject, OnInit } from 'final-di';
import { SearchClient, SearchClientResponse } from './search-client';
import { PostgreAdapterService } from '../postgre/postgre-adapter.service';
import { PostgreService } from '../postgre/postgre.service';
import { Logger } from '../../infrastructure/utils/logger';
import { TopicRepository } from '../repositories/topics/topic-repository';
import { SearchRepository } from '../repositories/search-repository';
import { MotionBlockRepository } from '../repositories/motions/motion-block-repository';

export class SearchService implements SearchClient, OnInit {
    @Inject(PostgreAdapterService)
    private readonly postgre!: PostgreService;

    @Inject(TopicRepository)
    private readonly topicRepo!: SearchRepository;

    @Inject(MotionBlockRepository)
    private readonly motionblockRepo!: SearchRepository;

    private readonly _repositories = [this.topicRepo, this.motionblockRepo];

    public async onInit(): Promise<void> {
        await this.initSearchIndexes();
        Logger.log(`Search indices created!`);
        await this.logCollections();
    }

    public async search(searchQuery: string): Promise<SearchClientResponse[]> {
        searchQuery = searchQuery.split(' ').join(' | ');
        Logger.debug(`Start search with query: ${searchQuery}`);
        const promises = this._repositories.map(repo => repo.search(searchQuery));
        const results = await Promise.all(promises);
        const logOutput = results.length === 1 ? `1 entry` : `${results.length} entries`;
        Logger.debug(`Search result contains ${logOutput}.`);
        for (const entry of results) {
            Logger.debug(JSON.stringify(entry));
        }
        return results;
    }

    private async initSearchIndexes(): Promise<void> {
        await Promise.all(this._repositories.map(repo => repo.buildColumn()));
        await Promise.all(this._repositories.map(repo => repo.buildIndex()));
        await Promise.all(this._repositories.map(repo => repo.buildTriggerFn()));
    }

    private async logCollections(): Promise<void> {
        const models = await this.postgre.getFqids();
        const collections: { [collection: string]: string[] } = {};
        for (const model of models) {
            const [collection, id] = model.fqid.split(`/`);
            if (!collections[collection]) {
                collections[collection] = [];
            }
            collections[collection].push(id);
        }
        Logger.debug(`Found ${Object.keys(collections).length} collections!`);
        for (const collection of Object.keys(collections)) {
            Logger.debug(`Found collection ${collection.toUpperCase()}: ${collections[collection].length} entries!`);
        }
    }
}
