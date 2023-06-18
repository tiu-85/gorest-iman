CREATE SCHEMA IF NOT EXISTS post AUTHORIZATION iman;

CREATE TABLE post.posts (
                            id serial4 NOT NULL,
                            post_id int4 NOT NULL,
                            user_id int4 NOT NULL,
                            title varchar NULL,
                            body varchar NULL,
                            CONSTRAINT posts_pk PRIMARY KEY (id)
);

CREATE TABLE post.tasks (
                            id serial4 NOT NULL,
                            total int4 NOT NULL DEFAULT 0,
                            success int4 NOT NULL DEFAULT 0,
                            fail int4 NOT NULL DEFAULT 0,
                            page_offset int4 NOT NULL DEFAULT 0,
                            page_limit int4 NOT NULL DEFAULT 0,
                            CONSTRAINT tasks_pk PRIMARY KEY (id)
);