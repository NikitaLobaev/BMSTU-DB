\c forums_3nf;

DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE EXTENSION citext;

CREATE TYPE voice AS ENUM ('1', '-1');

CREATE UNLOGGED TABLE user_ (
    nickname citext COLLATE "C" NOT NULL PRIMARY KEY,
    about TEXT NOT NULL DEFAULT '',
    email citext NOT NULL UNIQUE,
    fullname TEXT NOT NULL
);

CREATE UNLOGGED TABLE forum (
    slug citext NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    user_nickname citext NOT NULL REFERENCES user_ (nickname) ON DELETE CASCADE,
    threads INT NOT NULL DEFAULT 0,
    posts INT NOT NULL DEFAULT 0
);

CREATE UNLOGGED TABLE thread (
    id SERIAL PRIMARY KEY,
    user_nickname citext NOT NULL REFERENCES user_ (nickname) ON DELETE CASCADE,
    created TIMESTAMPTZ NOT NULL,
    forum_slug citext NOT NULL REFERENCES forum (slug) ON DELETE CASCADE,
    message TEXT NOT NULL,
    slug citext UNIQUE,
    title TEXT NOT NULL,
    votes INT NOT NULL DEFAULT 0
);

CREATE UNLOGGED TABLE post (
    id BIGSERIAL PRIMARY KEY,
    user_nickname citext NOT NULL REFERENCES user_ (nickname) ON DELETE CASCADE,
    created TIMESTAMP NOT NULL,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    message TEXT NOT NULL,
    post_parent_id BIGINT REFERENCES post (id) ON DELETE CASCADE,
    thread_id INT NOT NULL REFERENCES thread (id) ON DELETE CASCADE
);

CREATE UNLOGGED TABLE vote (
    user_nickname citext NOT NULL REFERENCES user_ (nickname) ON DELETE CASCADE,
    thread_id INT NOT NULL REFERENCES thread (id) ON DELETE CASCADE,
    PRIMARY KEY (user_nickname, thread_id),
    voice voice NOT NULL
);

CREATE UNLOGGED TABLE forum_user (
    forum_slug citext NOT NULL REFERENCES forum (slug) ON DELETE CASCADE,
    user_nickname citext COLLATE "C" NOT NULL REFERENCES user_ (nickname) ON DELETE CASCADE,
    PRIMARY KEY (forum_slug, user_nickname)
);

CREATE INDEX ON user_ USING hash (nickname);
CREATE INDEX ON user_ USING hash (email);

CREATE INDEX ON forum USING hash (slug);

CREATE INDEX ON thread USING hash (forum_slug);
CREATE INDEX ON thread (forum_slug, created);

CREATE INDEX ON post (thread_id);
CREATE INDEX ON post (created, id);

CREATE INDEX ON forum_user USING hash (forum_slug);

CREATE FUNCTION trigger_forum_before_insert()
    RETURNS TRIGGER
AS $$
BEGIN
    NEW.threads := 0;
    NEW.posts := 0;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert BEFORE INSERT
    ON forum
    FOR EACH ROW
EXECUTE PROCEDURE trigger_forum_before_insert();

CREATE FUNCTION trigger_thread_before_insert()
    RETURNS TRIGGER
AS $$
BEGIN
    IF NEW.slug = '' THEN
        NEW.slug := NULL;
    END IF;
    NEW.votes := 0;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert BEFORE INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE trigger_thread_before_insert();

CREATE FUNCTION add_forum_user(arg_forum_slug citext, arg_user_nickname citext)
    RETURNS VOID
AS $$
BEGIN
    INSERT INTO forum_user (forum_slug, user_nickname) VALUES (arg_forum_slug, arg_user_nickname)
    ON CONFLICT (forum_slug, user_nickname) DO NOTHING;
    RETURN;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION trigger_thread_after_insert()
    RETURNS TRIGGER
AS $$
BEGIN
    UPDATE forum SET threads = threads + 1 WHERE slug = NEW.forum_slug;
    EXECUTE add_forum_user(NEW.forum_slug, NEW.user_nickname);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE trigger_thread_after_insert();

CREATE FUNCTION trigger_post_before_insert()
    RETURNS TRIGGER
AS $$
BEGIN
    IF NEW.post_parent_id = 0 THEN
        NEW.post_parent_id := NULL;
    ELSEIF NEW.post_parent_id IS NOT NULL AND ((SELECT COUNT(*) FROM post WHERE id = NEW.post_parent_id) = 0 OR (SELECT thread_id FROM post WHERE id = NEW.post_parent_id) != NEW.thread_id) THEN
        RAISE 'Parent post is in another thread';
    END IF;
    NEW.is_edited := FALSE;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert BEFORE INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE trigger_post_before_insert();

CREATE FUNCTION trigger_post_after_insert()
    RETURNS TRIGGER
AS $$
    DECLARE fs citext;
BEGIN
    fs := (SELECT forum_slug FROM thread WHERE thread.id = NEW.thread_id);
    UPDATE forum SET posts = posts + 1 WHERE slug = fs;
    EXECUTE add_forum_user(fs, NEW.user_nickname);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert AFTER INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE trigger_post_after_insert();

CREATE FUNCTION trigger_post_before_update()
    RETURNS TRIGGER
AS $$
BEGIN
    NEW.is_edited := NEW.message != OLD.message;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update BEFORE UPDATE
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE trigger_post_before_update();

CREATE FUNCTION trigger_vote_after_insert()
    RETURNS TRIGGER
AS $$
BEGIN
    IF NEW.voice = '1' THEN
        UPDATE thread SET votes = votes + 1 WHERE id = NEW.thread_id;
    ELSE
        UPDATE thread SET votes = votes - 1 WHERE id = NEW.thread_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert AFTER INSERT
    ON vote
    FOR EACH ROW
EXECUTE PROCEDURE trigger_vote_after_insert();

CREATE FUNCTION trigger_vote_after_update()
    RETURNS TRIGGER
AS $$
BEGIN
    IF OLD.voice != NEW.voice THEN
        IF NEW.voice = '1' THEN
            UPDATE thread SET votes = votes + 2 WHERE id = NEW.thread_id;
        ELSE
            UPDATE thread SET votes = votes - 2 WHERE id = NEW.thread_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_update AFTER UPDATE
    ON vote
    FOR EACH ROW
EXECUTE PROCEDURE trigger_vote_after_update();

CREATE FUNCTION get_post_path(id_ BIGINT, OUT path_ BIGINT[])
AS $$
BEGIN
    path_ := ARRAY[id_];
    WHILE TRUE LOOP
        id_ := (SELECT post_parent_id FROM post WHERE id = id_);
        IF id_ IS NULL THEN
            EXIT;
        END IF;
        path_ := ARRAY[id_] || path_;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION get_post_root_id(id BIGINT, OUT root_id BIGINT)
AS $$
BEGIN
    root_id := (get_post_path(id))[1];
END;
$$ LANGUAGE plpgsql;
