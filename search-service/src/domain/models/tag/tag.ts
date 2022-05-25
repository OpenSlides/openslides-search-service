import { Fqid } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

export class Tag {
    public static COLLECTION = `tag`;

    public name!: string;

    public tagged_ids!: Fqid[]; // (*/tag_ids)[];
}
export interface Tag extends HasMeetingId {}
