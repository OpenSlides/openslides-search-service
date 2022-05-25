import { SearchRepository } from '../search-repository';
import { Topic } from '../../../domain/models/topics/topic';

export class TopicRepository extends SearchRepository<Topic> {
    public readonly COLLECTION = Topic.COLLECTION;

    public getSearchableFields(): (keyof Topic)[] {
        return [`text`, `title`];
    }
}
