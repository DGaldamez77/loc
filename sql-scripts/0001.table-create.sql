DROP SCHEMA IF EXISTS loc CASCADE;

CREATE SCHEMA loc;
ALTER SCHEMA loc OWNER TO loc_admin;

GRANT USAGE ON SCHEMA loc TO loc_admin;
GRANT USAGE ON SCHEMA loc TO loc_service;
GRANT USAGE ON SCHEMA loc TO loc_readonly;
GRANT USAGE ON SCHEMA loc TO loc_readwrite;


CREATE FUNCTION loc.table_changed() RETURNS TRIGGER
    LANGUAGE plpgsql
AS
$$
BEGIN
    IF NEW.* IS DISTINCT FROM OLD.* THEN
        NEW.updated := current_timestamp;
        IF NEW.UPDATED_BY IS NULL THEN
            NEW.updated_by := current_user;
        END IF;
    END IF;
    RETURN NEW;
END;
$$;


CREATE TABLE loc.genre
(
    genre_id            SERIAL
        CONSTRAINT genre_pk
            PRIMARY KEY,
    genre                    VARCHAR NOT NULL,
    active                  BOOLEAN NOT NULL         DEFAULT TRUE,
    created                 TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by              VARCHAR NOT NULL         DEFAULT CURRENT_USER,
    updated                 TIMESTAMP WITH TIME ZONE,
    updated_by              VARCHAR
);

ALTER TABLE loc.genre
    owner to loc_admin;

CREATE UNIQUE INDEX genre_genre_id_uindex 
    ON loc.genre (genre_id);

CREATE UNIQUE INDEX genre_genre_uindex 
    ON loc.genre (lower(genre));

GRANT ALL ON loc.genre TO loc_readwrite;
GRANT SELECT ON loc.genre TO loc_readonly;

CREATE TRIGGER genre_u
    BEFORE UPDATE
    ON loc.genre
    FOR EACH ROW
EXECUTE PROCEDURE loc.table_changed();



CREATE TABLE loc.publisher
(
    publisher_id            SERIAL
        CONSTRAINT publisher_pk
            PRIMARY KEY,
    publisher               VARCHAR NOT NULL,
    active                  BOOLEAN NOT NULL         DEFAULT TRUE,
    created                 TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by              VARCHAR NOT NULL         DEFAULT CURRENT_USER,
    updated                 TIMESTAMP WITH TIME ZONE,
    updated_by              VARCHAR
);

ALTER TABLE loc.publisher
    owner to loc_admin;

CREATE UNIQUE INDEX publisher_publisher_id_uindex 
    ON loc.publisher (publisher_id);

CREATE UNIQUE INDEX publisher_publisher_uindex 
    ON loc.publisher (lower(publisher));

GRANT ALL ON loc.publisher TO loc_readwrite;
GRANT SELECT ON loc.publisher TO loc_readonly;

CREATE TRIGGER publisher_u
    BEFORE UPDATE
    ON loc.publisher
    FOR EACH ROW
EXECUTE PROCEDURE loc.table_changed();


CREATE TABLE loc.book
(
    book_id            SERIAL
        CONSTRAINT book_pk
            PRIMARY KEY,
    title                   VARCHAR NOT NULL,
    publisher_id            INT
        CONSTRAINT book_publisher_id_fk
            REFERENCES loc.publisher,
    isbn                    VARCHAR NOT NULL,
    synopsis                VARCHAR NOT NULL,
    genre_id                INT
        CONSTRAINT book_genre_id_fk
            REFERENCES loc.genre,
    active                  BOOLEAN NOT NULL         DEFAULT TRUE,
    created                 TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by              VARCHAR NOT NULL         DEFAULT CURRENT_USER,
    updated                 TIMESTAMP WITH TIME ZONE,
    updated_by              VARCHAR
);

ALTER TABLE loc.book
    owner to loc_admin;

CREATE UNIQUE INDEX book_book_id_uindex 
    ON loc.book (book_id);

CREATE INDEX book_title_index 
    ON loc.book (lower(title));

GRANT ALL ON loc.book TO loc_readwrite;
GRANT SELECT ON loc.book TO loc_readonly;

CREATE TRIGGER book_u
    BEFORE UPDATE
    ON loc.book
    FOR EACH ROW
EXECUTE PROCEDURE loc.table_changed();



CREATE TABLE loc.author
(
    author_id            SERIAL
        CONSTRAINT author_pk
            PRIMARY KEY,
    first_name              VARCHAR NOT NULL,
    last_name               VARCHAR NOT NULL,
    active                  BOOLEAN NOT NULL         DEFAULT TRUE,
    created                 TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by              VARCHAR NOT NULL         DEFAULT CURRENT_USER,
    updated                 TIMESTAMP WITH TIME ZONE,
    updated_by              VARCHAR
);

ALTER TABLE loc.author
    owner to loc_admin;

CREATE UNIQUE INDEX author_author_id_uindex 
    ON loc.author (author_id);

CREATE UNIQUE INDEX author_first_name_last_name_uindex 
    ON loc.author (lower(first_name), lower(last_name));

GRANT ALL ON loc.author TO loc_readwrite;
GRANT SELECT ON loc.author TO loc_readonly;

CREATE TRIGGER author_u
    BEFORE UPDATE
    ON loc.author
    FOR EACH ROW
EXECUTE PROCEDURE loc.table_changed();



CREATE TABLE loc.book_author_xref
(
    book_author_xref_id            SERIAL
        CONSTRAINT book_author_pk
            PRIMARY KEY,
    book_id                 INT
        CONSTRAINT book_author_bool_id_fk
            REFERENCES loc.book,
    author_id               INT
        CONSTRAINT book_author_author_id_fk
            REFERENCES loc.author,
    created                 TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by              VARCHAR NOT NULL         DEFAULT CURRENT_USER,
    updated                 TIMESTAMP WITH TIME ZONE,
    updated_by              VARCHAR
);

ALTER TABLE loc.book_author_xref
    owner to loc_admin;

CREATE UNIQUE INDEX book_author_xref_book_author_xref_id_uindex 
    ON loc.book_author_xref (book_author_xref_id);

CREATE UNIQUE INDEX book_author_xref_book_id_uindex 
    ON loc.book_author_xref (book_id);

GRANT ALL ON loc.book_author_xref TO loc_readwrite;
GRANT SELECT ON loc.book_author_xref TO loc_readonly;

CREATE TRIGGER book_author_xref_u
    BEFORE UPDATE
    ON loc.book_author_xref
    FOR EACH ROW
EXECUTE PROCEDURE loc.table_changed();



CREATE TABLE loc.user
(
    user_id            SERIAL
        CONSTRAINT user_pk
            PRIMARY KEY,
    user_name               VARCHAR NOT NULL,
    email                   VARCHAR NOT NULL,
    active                  BOOLEAN NOT NULL         DEFAULT TRUE,
    created                 TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by              VARCHAR NOT NULL         DEFAULT CURRENT_USER,
    updated                 TIMESTAMP WITH TIME ZONE,
    updated_by              VARCHAR
);

ALTER TABLE loc.user
    owner to loc_admin;

CREATE UNIQUE INDEX user_user_id_uindex 
    ON loc.user (user_id);

CREATE UNIQUE INDEX user_name_uindex 
    ON loc.user (user_name);

CREATE UNIQUE INDEX email_uindex 
    ON loc.user (email);

GRANT ALL ON loc.user TO loc_readwrite;
GRANT SELECT ON loc.user TO loc_readonly;

CREATE TRIGGER user_u
    BEFORE UPDATE
    ON loc.user
    FOR EACH ROW
EXECUTE PROCEDURE loc.table_changed();


CREATE TABLE loc.book_inventory
(
    -- book_inventory_id            SERIAL
    --     CONSTRAINT book_inventory_pk
    --         PRIMARY KEY,
    book_id                 INT
        CONSTRAINT book_inventory_pk
            PRIMARY KEY
        CONSTRAINT book_inventory_book_id_fk
            REFERENCES loc.book,
    book_count              INT NOT NULL,
    created                 TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by              VARCHAR NOT NULL         DEFAULT CURRENT_USER,
    updated                 TIMESTAMP WITH TIME ZONE,
    updated_by              VARCHAR
);

ALTER TABLE loc.book_inventory
    owner to loc_admin;

CREATE UNIQUE INDEX book_inventory_book_id_uindex 
    ON loc.book_inventory (book_id);

GRANT ALL ON loc.book_inventory TO loc_readwrite;
GRANT SELECT ON loc.book_inventory TO loc_readonly;

CREATE TRIGGER book_inventory_u
    BEFORE UPDATE
    ON loc.book_inventory
    FOR EACH ROW
EXECUTE PROCEDURE loc.table_changed();



CREATE TABLE loc.book_checkout
(
    book_checkout_id            SERIAL
        CONSTRAINT book_checkout_pk
            PRIMARY KEY,
    book_id                 INT
        CONSTRAINT book_checkout_book_id_fk
            REFERENCES loc.book,
    user_id                 INT
        CONSTRAINT book_checkout_user_id_fk
            REFERENCES loc.user,
    returned                TIMESTAMP WITH TIME ZONE,
    created                 TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_by              VARCHAR NOT NULL         DEFAULT CURRENT_USER,
    updated                 TIMESTAMP WITH TIME ZONE,
    updated_by              VARCHAR
);

ALTER TABLE loc.book_checkout
    owner to loc_admin;

CREATE UNIQUE INDEX book_checkout_book_checkout_id_uindex 
    ON loc.book_checkout (book_checkout_id);

CREATE INDEX book_checkout_book_id_user_id_index 
    ON loc.book_checkout (book_id, user_id);

GRANT ALL ON loc.book_checkout TO loc_readwrite;
GRANT SELECT ON loc.book_checkout TO loc_readonly;

CREATE TRIGGER book_checkout_u
    BEFORE UPDATE
    ON loc.book_checkout
    FOR EACH ROW
EXECUTE PROCEDURE loc.table_changed();

GRANT ALL ON SEQUENCE loc.book_checkout_book_checkout_id_seq TO loc_readwrite;
GRANT ALL ON SEQUENCE loc.book_book_id_seq TO loc_readwrite;
GRANT ALL ON SEQUENCE loc.author_author_id_seq TO loc_readwrite;
GRANT ALL ON SEQUENCE loc.genre_genre_id_seq TO loc_readwrite;
GRANT ALL ON SEQUENCE loc.book_author_xref_book_author_xref_id_seq TO loc_readwrite;
GRANT ALL ON SEQUENCE loc.publisher_publisher_id_seq TO loc_readwrite;
GRANT ALL ON SEQUENCE loc.user_user_id_seq TO loc_readwrite;
