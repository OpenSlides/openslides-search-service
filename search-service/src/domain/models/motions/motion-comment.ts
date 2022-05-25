import { Id } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

export class MotionComment {
    public static COLLECTION = `motion_comment`;

    public comment!: string;

    public motion_id!: Id; // motion/comment_ids;
    public section_id!: Id; // motion_comment_section/comment_ids;
}
export interface MotionComment extends HasMeetingId {}
