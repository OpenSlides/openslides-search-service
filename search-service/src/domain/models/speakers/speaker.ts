import { Id, UnsafeHtml } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

/**
 * Representation of a speaker in a list of speakers.
 */
export class Speaker {
    public static COLLECTION = `speaker`;

    /**
     * Unixtime. Null if the speaker has not started yet. This time is in seconds.
     */
    public begin_time!: number;

    /**
     * Unixtime. Null if the speech has not ended yet. This time is in seconds.
     */
    public end_time!: number;

    public weight!: number;
    public point_of_order!: boolean;
    public speech_state!: number;
    public note!: UnsafeHtml;

    public list_of_speakers_id!: Id; // list_of_speakers/speaker_ids;
    public user_id!: Id; // user/speaker_$<meeting_id>_ids;

    public get speakingTime(): number {
        return this.end_time - this.begin_time || 0;
    }
}
export interface Speaker extends HasMeetingId {}
