import { SearchRepository } from '../search-repository';
import { MotionBlock } from '../../../domain/models/motions/motion-block';

export class MotionBlockRepository extends SearchRepository<MotionBlock> {
    public COLLECTION: string = MotionBlock.COLLECTION;

    public getSearchableFields(): (keyof MotionBlock)[] {
        return ['title', 'internal'];
    }
}
