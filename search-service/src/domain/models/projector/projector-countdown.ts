import { HasMeetingId } from '../../interfaces/has-meeting-id';
import { HasProjectionIds } from '../../interfaces/has-projectable-ids';

export class ProjectorCountdown {
    public static COLLECTION = `projector_countdown`;

    public title!: string;
    public description!: string;
    public default_time!: number;
    public countdown_time!: number;
    public running!: boolean;
}
export interface ProjectorCountdown extends HasMeetingId, HasProjectionIds {}
