import { Id } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

export class MotionSubmitter {
    public static COLLECTION = `motion_submitter`;

    public weight!: number;

    public user_id!: Id; // user/submitted_motion_$<meeting_id>_ids;
    public motion_id!: Id; // motion/submitter_ids;
}
export interface MotionSubmitter extends HasMeetingId {}
