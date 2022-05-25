import { Id } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

export class MotionStatuteParagraph {
    public static COLLECTION = `motion_statute_paragraph`;

    public title!: string;
    public text!: string;
    public weight!: number;

    public motion_ids!: Id[]; // (motion/statute_paragraph_id)[];
}
export interface MotionStatuteParagraph extends HasMeetingId {}
