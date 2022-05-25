import { Fqid, Id, Ids } from '../../definitions/key-types';

export class Option {
    public static COLLECTION = `option`;

    public text!: string;
    public yes!: number;
    public no!: number;
    public abstain!: number;

    public poll_id!: Id; // (assignment|motion)_poll/option_ids;
    public vote_ids!: Ids; // ((assignment|motion)_vote/option_id)[];

    public weight!: number;
    public content_object_id!: Fqid;
}
