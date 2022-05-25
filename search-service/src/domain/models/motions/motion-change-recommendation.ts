import { Id } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

export class MotionChangeRecommendation {
    public static COLLECTION = `motion_change_recommendation`;

    public rejected!: boolean;
    public internal!: boolean;
    public type!: string;
    public other_description!: string;
    public line_from!: number;
    public line_to!: number;
    public text!: string;
    public creation_time!: string;

    public motion_id!: Id; // motion/change_recommendation_ids;
}
export interface MotionChangeRecommendation extends HasMeetingId {}
