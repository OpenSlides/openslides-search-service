import { Id } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

export class MotionState {
    public static COLLECTION = `motion_state`;

    public name!: string;
    public recommendation_label!: string;
    public css_class!: string;
    public restrictions!: string[];
    public allow_support!: boolean;
    public allow_create_poll!: boolean;
    public allow_submitter_edit!: boolean;
    public allow_motion_forwarding!: boolean;
    public set_created_timestamp!: boolean;
    public set_number!: boolean;
    public show_state_extension_field!: boolean;
    public merge_amendment_into_final!: string;
    public show_recommendation_extension_field!: boolean;
    public weight!: number;

    public next_state_ids!: Id[]; // (motion_state/previous_state_ids)[];
    public previous_state_ids!: Id[]; // (motion_state/next_state_ids)[];
    public motion_ids!: Id[]; // (motion/state_id)[];
    public motion_recommendation_ids!: Id[]; // (motion/recommendation_id)[];
    public workflow_id!: Id; // motion_workflow/state_ids;
    public first_state_of_workflow_id!: Id; // motion_workflow/first_state_id;
}
export interface MotionState extends HasMeetingId {}
