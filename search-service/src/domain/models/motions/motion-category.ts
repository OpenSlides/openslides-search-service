import { Id } from '../../definitions/key-types';
import { HasMeetingId, HasSequentialNumber } from '../../interfaces';

export class MotionCategory {
    public static COLLECTION = `motion_category`;

    public name!: string;
    public prefix!: string;
    public weight!: number;
    public level!: number;

    public parent_id!: Id; // motion_category/child_ids;
    public child_ids!: Id[]; // (motion_category/parent_id)[];
    public motion_ids!: Id[]; // (motion/category_id)[];
}
export interface MotionCategory extends HasMeetingId, HasSequentialNumber {}
