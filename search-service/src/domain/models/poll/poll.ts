import { HasSequentialNumber } from '../../interfaces';

import { Fqid, Id } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';
import { HasProjectionIds } from '../../interfaces/has-projectable-ids';

export class Poll {
    public static readonly COLLECTION = `poll`;
    public static readonly DECIMAL_FIELDS: (keyof Poll)[] = [`votesvalid`, `votesinvalid`, `votescast`];

    public content_object_id!: Fqid;
    public state!: string;
    public type!: string;
    public title!: string;
    public votesvalid!: number;
    public votesinvalid!: number;
    public votescast!: number;
    public vote_count!: number;
    public onehundred_percent_base!: string;

    /**
     * TODO:
     * Not sure how vote delegations are handled now
     */
    public user_has_voted_for_delegations!: Id[];

    public pollmethod!: string;

    public voted_ids!: Id[]; // (user/poll_voted_$<meeting_id>_ids)[];

    public entitled_group_ids!: Id[]; // (group/(assignment|motion)_poll_ids)[];
    public option_ids!: Id[]; // ((assignment|motion)_option/poll_id)[];
    public global_option_id!: Id; // (motion_option/poll_id)
    public backend!: string;

    public description!: string;
    public min_votes_amount!: number;
    public max_votes_amount!: number;
    public max_votes_per_option!: number;
    public global_yes!: boolean;
    public global_no!: boolean;
    public global_abstain!: boolean;
    public entitled_users_at_stop!: Id[];
}

export interface Poll extends HasMeetingId, HasProjectionIds, HasSequentialNumber {}
