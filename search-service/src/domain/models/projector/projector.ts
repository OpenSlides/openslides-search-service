import { Id } from '../../definitions/key-types';
import { HasSequentialNumber } from '../../interfaces';
import { HasMeetingId } from '../../interfaces/has-meeting-id';
import { Projectiondefault } from './projection-default';

/**
 * Representation of a projector.
 */
export class Projector {
    public static COLLECTION = `projector`;

    public name!: string;
    public scale!: number;
    public scroll!: number;
    public width!: number;
    public aspect_ratio_numerator!: number;
    public aspect_ratio_denominator!: number;
    public color!: string;
    public background_color!: string;
    public header_background_color!: string;
    public header_font_color!: string;
    public header_h1_color!: string;
    public chyron_background_color!: string;
    public chyron_font_color!: string;
    public show_header_footer!: boolean;
    public show_title!: boolean;
    public show_logo!: boolean;
    public show_clock!: boolean;

    public current_projection_ids!: Id[]; // (projection/current_projector_id)[];
    public preview_projection_ids!: Id[]; // (projection/preview_projector_id)[];
    public history_projection_ids!: Id[]; // (projection/history_projector_id)[];
    public used_as_reference_projector_meeting_id!: Id; // meeting/reference_projector_id;
    public used_as_default_$_in_meeting_id!: Projectiondefault[]; // meeting/default_projector_$_id;
}
export interface Projector extends HasMeetingId, HasSequentialNumber {}
