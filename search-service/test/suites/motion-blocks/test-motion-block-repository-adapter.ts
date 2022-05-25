import { MotionBlock } from '../../../src/domain/models/motions/motion-block';
import { Fqid } from '../../../src/domain/definitions/key-types';
import { RepositoryAdapter } from '../../utils/repository-adapter';

export class TestMotionBlockRepositoryAdapter extends RepositoryAdapter<MotionBlock> {
    protected getFqid(id: number): Fqid {
        return `${MotionBlock.COLLECTION}/${id}`;
    }
}
