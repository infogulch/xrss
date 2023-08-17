CREATE TABLE IF NOT EXISTS feed (
  id INTEGER PRIMARY KEY,
  url TEXT,
  title TEXT
)
STRICT;

CREATE TABLE IF NOT EXISTS fetch (
  id INTEGER PRIMARY KEY,
  feed_id INTEGER NOT NULL,
  added_count INTEGER NOT NULL,
  data TEXT CHECK (JSON_VALID(data)),
  FOREIGN KEY (feed_id) REFERENCES feed (id)
)
STRICT;

DROP VIEW IF EXISTS v_fetch;
CREATE VIEW v_fetch AS
  SELECT
    id,
    feed_id,
    data,
    added_count,
    data ->> '$.fetched_on' AS fetched_on,
    data ->> '$.fetched_duration' AS fetched_duration,
    data ->> '$.fetched_bytes' AS fetched_bytes,
    data ->> '$.fetched_count' AS fetched_count
  FROM fetch;

DROP VIEW IF EXISTS v_feed_fetch_info;
CREATE VIEW v_feed_fetch_info AS
  SELECT
    feed.id AS feed_id,
    url,
    last_modified,
    etag
  FROM
    feed
    LEFT JOIN (
      SELECT
        feed_id,
        data ->> "$.last_modified" AS last_modified,
        data ->> "$.etag" AS etag
      FROM fetch
      GROUP BY feed_id
      HAVING id = MAX(id)
    ) AS latest_fetch
      ON feed_id = feed.id;

CREATE TABLE IF NOT EXISTS item (
  id INTEGER PRIMARY KEY,
  feed_id INTEGER,
  fetch_id INTEGER,
  data TEXT CHECK (JSON_VALID(data)),
  guid TEXT GENERATED ALWAYS AS (data ->> '$.guid') STORED,
  published TEXT GENERATED ALWAYS AS (data ->> '$.publishedParsed'),
  FOREIGN KEY (feed_id) REFERENCES feed (id),
  FOREIGN KEY (fetch_id) REFERENCES fetch (id)
)
STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS item_ids ON item (feed_id, guid, fetch_id);

DROP VIEW IF EXISTS v_item;
CREATE VIEW v_item AS
  SELECT
    id,
    feed_id,
    fetch_id,
    data,
    guid,
    data ->> '$.title' AS title,
    data ->> '$.description' AS description,
    data ->> '$.publishedParsed' AS published,
    data ->> '$.link' AS link,
    data ->> '$.image' AS image,
    data ->> '$.authors' AS authors,
    data ->> '$.authors[0].name' AS author
  FROM item;

CREATE TABLE IF NOT EXISTS user (
  id INTEGER PRIMARY KEY,
  username TEXT UNIQUE NOT NULL
)
STRICT;

CREATE TABLE IF NOT EXISTS subscription (
  user_id INTEGER,
  feed_id INTEGER,
  FOREIGN KEY (user_id) REFERENCES user (id),
  FOREIGN KEY (feed_id) REFERENCES feed (id),
  PRIMARY KEY (user_id, feed_id)
)
STRICT;

CREATE TABLE IF NOT EXISTS seen (
  user_id INTEGER NOT NULL,
  feed_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES user (id),
  FOREIGN KEY (item_id) REFERENCES item (id),
  FOREIGN KEY (feed_id) REFERENCES feed (id),
  PRIMARY KEY (user_id, feed_id, item_id),
  UNIQUE (user_id, item_id)
)
STRICT;

DROP VIEW IF EXISTS v_unseen;
CREATE VIEW v_unseen AS
  SELECT
    subscription.user_id,
    subscription.feed_id,
    item.id AS item_id
  FROM
    subscription
    JOIN item ON subscription.feed_id = item.feed_id
  EXCEPT
  SELECT
    user_id,
    feed_id,
    item_id
  FROM seen;

DROP VIEW IF EXISTS v_unseen_counts;
CREATE VIEW v_unseen_counts AS
  SELECT
    user_id,
    feed_id,
    COUNT(item_id) AS unread_count
  FROM v_unseen
  GROUP BY user_id, feed_id;

DROP VIEW IF EXISTS v_feeds_list;
CREATE VIEW v_feeds_list AS
  SELECT
    user_id,
    feed_id,
    title,
    unread_count
  FROM
    v_unseen_counts
    JOIN feed ON feed.id = feed_id;

DROP VIEW IF EXISTS ingest;
CREATE VIEW ingest AS
  SELECT
    fetch.id AS fetch_id,
    fetch.feed_id,
    added_count,
    JSON_SET(fetch.data, "$.feed.items", items.items) AS data
  FROM
    fetch
    JOIN feed ON fetch.feed_id = feed.id
    JOIN (
      SELECT
        feed_id,
        fetch_id,
        JSON_GROUP_ARRAY(data) AS items
      FROM item
      GROUP BY feed_id, fetch_id
    ) AS items
      ON items.fetch_id = fetch.id AND items.feed_id = fetch.feed_id;

CREATE TRIGGER ingest_fetch
INSTEAD OF INSERT ON ingest
FOR EACH ROW
BEGIN
  INSERT INTO fetch
    (feed_id, data, added_count)
  VALUES
    (NEW.feed_id, JSON_REMOVE(NEW.data, "$.feed.items"), 0);

  INSERT INTO item
    (feed_id, fetch_id, data)
  SELECT
    NEW.feed_id,
    (SELECT MAX(id) FROM fetch) AS fetch_id,
    JSON_EACH.value AS data
  FROM JSON_EACH(NEW.data, "$.feed.items")
  WHERE
    JSON_EACH.value ->> "$.guid" NOT IN (
      SELECT guid FROM item WHERE feed_id = NEW.feed_id
    );

  UPDATE fetch
  SET added_count = added_count + CHANGES()
  WHERE id = (SELECT MAX(id) FROM fetch);

  UPDATE feed
  SET title = NEW.data ->> "$.feed.title"
  WHERE feed.id = NEW.feed_id AND title IS NULL;
END;
