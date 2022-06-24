CREATE EXTENSION IF NOT EXISTS CITEXT;
DROP TABLE IF EXISTS Users CASCADE;
DROP TABLE IF EXISTS Forums CASCADE;
DROP TABLE IF EXISTS Threads CASCADE;
DROP TABLE IF EXISTS Posts CASCADE;
DROP TABLE IF EXISTS ForumUsers CASCADE;
DROP TABLE IF EXISTS Votes CASCADE;


CREATE UNLOGGED TABLE IF NOT EXISTS Users (
    id SERIAL,
    nickname CITEXT COLLATE "C" NOT NULL PRIMARY KEY,
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT NOT NULL UNIQUE
);

DROP INDEX user_nickname;
CREATE INDEX IF NOT EXISTS user_nickname ON Users(nickname);

DROP INDEX user_email;
CREATE INDEX IF NOT EXISTS user_email ON Users USING hash(email);

CREATE UNLOGGED TABLE IF NOT EXISTS Forums (
    id  SERIAL,
    slug CITEXT PRIMARY KEY,
    title TEXT NOT NULL,
    "user" CITEXT COLLATE "C" NOT NULL REFERENCES Users(nickname),
    posts INT NOT NULL DEFAULT 0,
    threads INT NOT NULL DEFAULT 0
);

DROP INDEX forum_slug;
CREATE INDEX IF NOT EXISTS forum_slug ON Forums using hash(slug);

CREATE UNLOGGED TABLE Threads (
    id SERIAL NOT NULL PRIMARY KEY,
    slug CITEXT UNIQUE,
    title TEXT NOT NULL,
    author CITEXT COLLATE "C" NOT NULL REFERENCES Users(nickname),
    forum CITEXT NOT NULL REFERENCES Forums(slug) ,
    message TEXT,
    votes INT DEFAULT 0,
    created TIMESTAMP WITH TIME ZONE DEFAULT now()
);

DROP INDEX thread_slug;
CREATE INDEX IF NOT EXISTS thread_slug ON Threads USING hash(slug);

DROP INDEX thread_forum_slug;
CREATE INDEX IF NOT EXISTS thread_forum_slug ON Threads(forum);

DROP INDEX thread_forum_created_idx;
CREATE INDEX IF NOT EXISTS thread_forum_created_idx ON Threads(slug, created);

CREATE UNLOGGED TABLE IF NOT EXISTS Posts (
    id SERIAL PRIMARY KEY,
    parent INT DEFAULT 0,
    path INT[] DEFAULT ARRAY []::INT[]
    author CITEXT  COLLATE "C" NOT NULL REFERENCES Users(nickname),
    message TEXT NOT NULL,
    isEdited boolean DEFAULT FALSE,
    forum CITEXT NOT NULL REFERENCES Forums(slug),
    thread INT  REFERENCES Threads(id),
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
);

DROP INDEX post_thread;
CREATE INDEX IF NOT EXISTS post_thread ON Posts(thread);

DROP INDEX post_thread_created;
CREATE INDEX IF NOT EXISTS post_thread_created ON Posts(thread_id,created);

DROP INDEX post_path;
CREATE INDEX IF NOT EXISTS post_path ON Posts((path[1]), path);

DROP INDEX post_thread_path;
CREATE INDEX IF NOT EXISTS post_thread_path ON Posts(thread, path);

CREATE UNLOGGED TABLE IF NOT EXISTS ForumUsers (
    nickname CITEXT COLLATE "C" NOT NULL REFERENCES Users(nickname),
    fullname TEXT NOT NULL,
    about TEXT,
    email CITEXT NOT NULL,
    forum CITEXT NOT NULL REFERENCES Forums(slug),
    PRIMARY KEY (nickname, forum)
);

DROP INDEX forum_users_forum_nickname;
CREATE INDEX IF NOT EXISTS forum_users_forum_nickname ON ForumUsers (forum, nickname);

CREATE UNLOGGED TABLE if not exists Votes (
    nickname CITEXT COLLATE "C" NOT NULL REFERENCES Users (nickname),
    thread SERIAL NOT NULL REFERENCES Threads(id),
    voice INT NOT NULL,
    PRIMARY KEY (nickname, thread)
);

DROP INDEX votes_nickname_thread;
CREATE UNIQUE INDEX IF NOT EXISTS votes_nickname_thread ON Votes (nickname, thread);

CREATE OR REPLACE FUNCTION insert_vote() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE Threads SET votes = votes + NEW.vote WHERE id = NEW.thread;
        RETURN NEW;
    END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_vote() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE Threads SET votes = votes - OLD.vote + NEW.vote WHERE id = NEW.thread;
        RETURN NEW;
    END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION increase_count_forum_threads() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE Forums SET threads = Forums.threads + 1 WHERE slug = NEW.forum;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION increase_count_forum_posts() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE Forums SET posts = Forums.posts + 1 WHERE slug = NEW.forum;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_post_path() RETURNS TRIGGER AS $$
    BEGIN
        new.path = (SELECT path FROM Posts WHERE id = new.parent) || new.id;
        RETURN new;
    END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION update_users_from_forum() RETURNS TRIGGER AS $$
DECLARE
nickname citext;
    fullname TEXT;
    about    TEXT;
    email    CITEXT;
    BEGIN
        SELECT u.nickname, u.fullname, u.about, u.email FROM users u WHERE u.nickname = NEW.author INTO nickname, fullname, about, email;
        INSERT INTO ForumUsers (nickname, fullname, about, email, forum)
        VALUES (nickname, fullname, about, email, NEW.forum) ON CONFLICT do nothing;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS vote_insert ON Votes;
CREATE TRIGGER insert_votes AFTER INSERT ON Votes FOR EACH ROW EXECUTE PROCEDURE insert_vote();

DROP TRIGGER IF EXISTS vote_insert ON Votes;
CREATE TRIGGER update_votes AFTER UPDATE ON Votes FOR EACH ROW EXECUTE PROCEDURE update_vote();

DROP TRIGGER IF EXISTS count_threads ON Threads;
CREATE TRIGGER count_threads AFTER INSERT ON Threads FOR EACH ROW EXECUTE PROCEDURE increase_count_forum_threads();

DROP TRIGGER IF EXISTS count_posts ON Posts;
CREATE TRIGGER count_posts AFTER INSERT ON Posts FOR EACH ROW EXECUTE PROCEDURE increase_count_forum_posts();

DROP TRIGGER IF EXISTS count_posts ON Posts;
CREATE TRIGGER update_post_path BEFORE INSERT ON Posts FOR EACH ROW EXECUTE PROCEDURE update_post_path();

DROP TRIGGER IF EXISTS update_users_on_post ON Posts;
CREATE TRIGGER update_users_on_post AFTER INSERT ON Posts FOR EACH ROW EXECUTE PROCEDURE update_users_from_forum();

DROP TRIGGER IF EXISTS update_users_on_thread ON Threads;
CREATE TRIGGER update_users_on_thread AFTER INSERT ON Threads FOR EACH ROW EXECUTE PROCEDURE update_users_from_forum();

VACUUM;
