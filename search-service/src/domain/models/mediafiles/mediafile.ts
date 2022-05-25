import { Fqid, Id } from '../../definitions/key-types';
import { HasListOfSpeakersId } from '../../interfaces/has-list-of-speakers-id';
import { HasOwnerId } from '../../interfaces/has-owner-id';
import { HasProjectionIds } from '../../interfaces/has-projectable-ids';

export class Mediafile {
    public static COLLECTION = `mediafile`;
    public static MEDIA_URL_PREFIX = `/media/`;

    public title!: string;
    public is_directory!: boolean;
    public filesize!: string;
    public filename!: string;
    public mimetype!: string;
    public pdf_information!: any;
    public create_timestamp!: string;
    public has_inherited_access_groups!: boolean;

    public access_group_ids!: Id[]; // (group/mediafile_access_group_ids)[];
    public inherited_access_group_ids!: Id[]; // (group/mediafile_inherited_access_group_ids)[];  // Note: calculated
    public parent_id!: Id; // mediafile/child_ids;
    public child_ids!: Id[]; // (mediafile/parent_id)[];
    public attachment_ids!: Fqid[]; // (*/attachment_ids)[];
    public used_as_logo_$_in_meeting_id!: string[]; // meeting/logo_$<place>_id;
    public used_as_font_$_in_meeting_id!: string[]; // meeting/font_$<place>_id;
}
export interface Mediafile extends HasOwnerId, HasProjectionIds, HasListOfSpeakersId {}
