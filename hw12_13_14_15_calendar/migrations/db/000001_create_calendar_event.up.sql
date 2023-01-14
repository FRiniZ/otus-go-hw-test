BEGIN;

CREATE TABLE IF NOT EXISTS events(
   id               SERIAL PRIMARY KEY,
   title            VARCHAR (150) NOT NULL,
   ontime           TIMESTAMP NOT NULL,
   offtime          TIMESTAMP,
   Description      TEXT,
   userid           BIGINT NOT NULL,
   NotifyTime       TIMESTAMP
);

CREATE INDEX IF NOT EXISTS events_userid_idx ON events (userid);
CREATE INDEX IF NOT EXISTS events_ontime_idx ON events (ontime);

COMMIT;