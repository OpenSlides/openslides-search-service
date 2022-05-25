import { Id } from '../../definitions/key-types';
import { HasProjectionIds } from '../../interfaces/has-projectable-ids';
import { Identifiable } from '../../interfaces/identifiable';

/**
 * Representation of a user in contrast to the operator.
 */
export class User {
    public static COLLECTION = `user`;

    public readonly username!: string;
    public readonly title!: string;
    public readonly pronoun!: string;
    public readonly first_name!: string;
    public readonly last_name!: string;
    public readonly is_active!: boolean;
    public readonly is_physical_person!: boolean;
    public readonly default_password!: string;
    public readonly can_change_own_password!: boolean;
    public readonly gender!: string;
    public readonly comment_$!: string[];
    public readonly number_$!: string[];
    public readonly about_me_$!: string[];
    public readonly default_number!: string;
    public readonly default_structure_level!: string;
    public readonly structure_level_$!: string[];
    public readonly email!: string;
    public readonly last_email_send!: number; // comes in seconds
    public readonly vote_weight_$!: number[];
    public readonly default_vote_weight!: number;
    public readonly is_demo_user!: boolean;

    // Meeting and committee
    public meeting_ids!: Id[]; // (meeting/user_ids)[];
    public is_present_in_meeting_ids!: Id[]; // (meeting/present_user_ids)[];
    public committee_ids!: Id[]; // (committee/user_ids)[];

    public group_$_ids!: string[]; // (group/user_ids)[];
    public speaker_$_ids!: string[]; // (speaker/user_id)[];
    public personal_note_$_ids!: string[]; // (personal_note/user_id)[];
    public supported_motion_$_ids!: string[]; // (motion/supporter_ids)[];
    public submitted_motion_$_ids!: string[]; // (motion_submitter/user_id)[];
    public poll_voted_$_ids!: string[]; // (poll/voted_ids)[];
    public vote_$_ids!: string[]; // (vote/user_id)[];
    public delegated_vote_$_ids!: string[]; // (vote/delegated_user_id)[];
    public option_$_ids!: string[];
    public assignment_candidate_$_ids!: string[]; // (assignment_candidate/user_id)[];
    public vote_delegated_$_vote_ids!: string[];
    public vote_delegated_$_to_id!: string[]; // user/vote_delegated_$<meeting_id>_from_ids;
    public vote_delegations_$_from_ids!: string[]; // user/vote_delegated_$<meeting_id>_to_id;
    public chat_message_$_ids!: Id[]; // (chat_message/user_id)[];

    public projection_$_ids!: any[];
    public current_projector_$_ids!: any[];

    public organization_management_level!: string;
    public committee_$_management_level!: string[];
}
export interface User extends Identifiable, HasProjectionIds {}
