import { Topic } from '../../../src/domain/models/topics/topic';
import { Fqid } from '../../../src/domain/definitions/key-types';
import { RepositoryAdapter } from '../../utils/repository-adapter';

export class TestTopicRepositoryAdapter extends RepositoryAdapter<Topic> {
    protected getFqid(id: number): Fqid {
        return `${Topic.COLLECTION}/${id}`;
    }
}
