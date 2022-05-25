import { Fqid, Id } from '../../definitions/key-types';
import { HasSequentialNumber } from '../../interfaces';
import { HasMeetingId } from '../../interfaces/has-meeting-id';
import { HasProjectionIds } from '../../interfaces/has-projectable-ids';

export class ListOfSpeakers {
    public static COLLECTION = `list_of_speakers`;

    public closed!: boolean;

    public content_object_id!: Fqid; // */list_of_speakers_id;
    public speaker_ids!: Id[]; // (speaker/list_of_speakers_id)[];
}
export interface ListOfSpeakers extends HasMeetingId, HasProjectionIds, HasSequentialNumber {}
