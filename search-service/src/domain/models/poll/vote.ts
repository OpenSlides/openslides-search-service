import { Id } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

export class Vote {
    public static COLLECTION = `vote`;

    public weight!: number;
    public value!: any;

    public option_id!: Id; // (assignment|motion)_option/vote_ids;
    public user_id!: Id; // user/(assignment|motion)_vote_$<meeting_id>_ids;
    public delegated_user_id!: Id; // user/(assignment|motion)_delegated_vote_$_ids;
    public user_token!: string;
}

export interface Vote extends HasMeetingId {}
