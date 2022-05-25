import { Id } from '../../definitions/key-types';
import { HasSequentialNumber } from '../../interfaces';
import { HasAgendaItemId } from '../../interfaces/has-agenda-item-id';
import { HasListOfSpeakersId } from '../../interfaces/has-list-of-speakers-id';
import { HasMeetingId } from '../../interfaces/has-meeting-id';
import { HasProjectionIds } from '../../interfaces/has-projectable-ids';
import { Identifiable } from '../../interfaces/identifiable';

export class MotionBlock {
    public static COLLECTION = `motion_block`;

    public title!: string;
    public internal!: boolean;

    public motion_ids!: Id[]; // (motion/block_id)[];
}
export interface MotionBlock
    extends Identifiable,
        HasMeetingId,
        HasAgendaItemId,
        HasListOfSpeakersId,
        HasProjectionIds,
        HasSequentialNumber {}
