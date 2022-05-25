import { HasMeetingId } from '../../interfaces/has-meeting-id';
import { HasProjectionIds } from '../../interfaces/has-projectable-ids';

export class ProjectorMessage {
    public static COLLECTION = `projector_message`;

    public message!: string;
}
export interface ProjectorMessage extends HasMeetingId, HasProjectionIds {}
