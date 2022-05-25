import { Id, UnsafeHtml } from '../../definitions/key-types';
import { HasMeetingId } from '../../interfaces/has-meeting-id';

export class ChatMessage {
    public static readonly COLLECTION = `chat_message`;

    public readonly meeting_id!: Id;

    public readonly created!: number; // in seconds
    public readonly content!: UnsafeHtml;

    public readonly user_id!: Id; // (user/chat_message_ids)
    public readonly chat_group_id!: Id; // (chat_group/chat_message_ids)
}

export interface ChatMessage extends HasMeetingId {}
