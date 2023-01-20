BEGIN;

DROP INDEX IF EXISTS events_userid_idx;
DROP INDEX IF EXISTS events_ontime_idx;
DROP INDEX IF EXISTS events_notify_idx;
DROP TABLE IF EXISTS events;

COMMIT;