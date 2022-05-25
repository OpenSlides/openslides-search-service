import { HasSequentialNumber } from '../../interfaces';
import { HasAgendaItemId } from '../../interfaces/has-agenda-item-id';
import { HasAttachmentIds } from '../../interfaces/has-attachment-ids';
import { HasListOfSpeakersId } from '../../interfaces/has-list-of-speakers-id';
import { HasMeetingId } from '../../interfaces/has-meeting-id';
import { HasTagIds } from '../../interfaces/has-tag-ids';
import { Identifiable } from '../../interfaces/identifiable';

export class Topic {
    public static COLLECTION = `topic`;

    public readonly title!: string;
    public readonly text!: string;
}

export interface Topic
    extends Identifiable,
        HasMeetingId,
        HasAgendaItemId,
        HasListOfSpeakersId,
        HasAttachmentIds,
        HasTagIds,
        HasSequentialNumber {}
