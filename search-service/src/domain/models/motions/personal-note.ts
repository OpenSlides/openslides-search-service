import { Fqid, Id } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

export class PersonalNote {
    public static COLLECTION = `personal_note`;

    public note!: string;
    public star!: boolean;

    public user_id!: Id; // user/personal_note_$<meeting_id>_ids;
    public content_object_id!: Fqid; // */personal_note_ids;
}
export interface PersonalNote extends HasMeetingId {}
