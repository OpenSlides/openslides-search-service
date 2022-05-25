import { Id } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

/**
 * Content of the 'assignment_related_users' property.
 */
export class AssignmentCandidate {
    public static COLLECTION = `assignment_candidate`;

    public weight!: number;

    public assignment_id!: Id; // assignment/candidate_ids;
    public user_id!: Id; // user/assignment_candidate_$<meeting_id>_ids;
}
export interface AssignmentCandidate extends HasMeetingId {}
