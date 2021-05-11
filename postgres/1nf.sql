\c forums_1nf;

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
    post_root_id BIGINT NOT NULL REFERENCES post (id) ON DELETE CASCADE,
    post_parent_id BIGINT REFERENCES post (id) ON DELETE CASCADE,
    path_ BIGINT[] NOT NULL,
    thread_id INT NOT NULL REFERENCES thread (id) ON DELETE CASCADE,
    forum_slug citext NOT NULL REFERENCES forum (slug) ON DELETE CASCADE
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
    PRIMARY KEY (forum_slug, user_nickname),
    user_about TEXT NOT NULL,
    user_email citext NOT NULL,
    user_fullname TEXT NOT NULL
);

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

CREATE FUNCTION trigger_user_after_update()
    RETURNS TRIGGER
AS $$
BEGIN
    IF OLD.about != NEW.about OR OLD.email != NEW.email OR OLD.fullname != NEW.fullname THEN
        UPDATE forum_user SET user_about = NEW.about, user_email = NEW.email, user_fullname = NEW.fullname
        WHERE forum_user.user_nickname = NEW.nickname;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_update AFTER INSERT
    ON user_
    FOR EACH ROW
EXECUTE PROCEDURE trigger_user_after_update();

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
    INSERT INTO forum_user (forum_slug, user_nickname, user_about, user_email, user_fullname)
    SELECT arg_forum_slug, arg_user_nickname, about, email, fullname FROM user_
    WHERE nickname = arg_user_nickname
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
    IF NEW.post_parent_id != 0 THEN
        NEW.path_ := (SELECT path_ FROM post WHERE thread_id = NEW.thread_id
                                               AND id = NEW.post_parent_id) || ARRAY[NEW.id];
        IF cardinality(NEW.path_) = 1 THEN
            RAISE 'Parent post is in another thread';
        END IF;
        NEW.post_root_id := NEW.path_[1];
    ELSE
        NEW.post_parent_id := NULL;
        NEW.post_root_id := NEW.id;
        NEW.path_ := ARRAY[NEW.id];
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
BEGIN
    UPDATE forum SET posts = posts + 1 WHERE slug = NEW.forum_slug;
    EXECUTE add_forum_user(NEW.forum_slug, NEW.user_nickname);
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
