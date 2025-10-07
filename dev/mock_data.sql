-- This script should only be used for an database filled with base_data.sql.

BEGIN;
INSERT INTO topic_t (id, title, text, sequential_number, meeting_id)
VALUES (2, 'Test Title', 'West test Wesley Teams and ', 1, 2);
SELECT nextval('topic_t_id_seq');

--list_of_speakers.content_object_id:topic.list_of_speakers_id gr:r
INSERT INTO list_of_speakers_t (
    id, content_object_id, sequential_number, meeting_id
)
VALUES (2, 'topic/2', 1, 2);

--agenda_item.content_object_id:topic.agenda_item_id gr:r
INSERT INTO agenda_item_t (content_object_id, meeting_id)
VALUES ('topic/2', 2);
COMMIT;

--rl:gr topic.poll_ids:poll.content_object_id
INSERT INTO poll_t (
    id,
    title,
    type,
    backend,
    pollmethod,
    onehundred_percent_base,
    sequential_number,
    content_object_id,
    meeting_id
)
VALUES (2, 'Titel1', 'analog', 'fast', 'YNA', 'disabled', 1, 'topic/2', 2);
SELECT nextval('poll_t_id_seq');

INSERT INTO gender_t (id, name) VALUES (2, 'test');

COMMIT;
