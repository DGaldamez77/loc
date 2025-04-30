TRUNCATE TABLE loc.book CASCADE;
TRUNCATE TABLE loc.author CASCADE;
TRUNCATE TABLE loc.book_author_xref CASCADE;
TRUNCATE TABLE loc.genre CASCADE;
TRUNCATE TABLE loc.publisher CASCADE;   
TRUNCATE TABLE loc.user CASCADE;
TRUNCATE TABLE loc.book_inventory CASCADE;

INSERT INTO loc.genre
	(genre_id, genre, active, created_by)
VALUES 
	(1, 'drama', true, 'david'),
	(2, 'fiction', true, 'david');

INSERT INTO loc.author 
	(author_id, first_name, last_name, active, created_by)
VALUES
	(1, 'John', 'Grisham', true, 'david'),
	(2, 'Mario', 'Puzo', true, 'david'),
	(3, 'Stephen', 'King', true, 'david'),
	(4, 'Dan', 'Brown', true, 'david');

INSERT INTO loc.publisher
	(publisher_id, publisher, active, created_by)
VALUES
	(1, 'Penguin Random House', true, 'david'),
	(2, 'HarperCollins', true, 'david');


INSERT INTO loc.book
	(book_id, title, publisher_id, isbn, synopsis, genre_id, active, created_by)
VALUES
	(1, 'The Godfather', 1, '978-0-439-02348-1','Best book ever', 1, true, 'david'),
	(2, 'The Sicilian', 1, '978-1-56619-909-4', '', 1, true, 'david'),
	(3, 'The Last Don', 1, '978-0-123456-47-2', '', 1, true, 'david'),
	(4, 'A Time To Kill', 2, '978-1-60309-502-0', '', 1, true, 'david'),
	(5, 'The Firm', 2, '978-1-60309-517-4', '', 1, true, 'david'),
	(6, 'The Pelican Brief', 2, '978-1-60309-442-9', '', 1, true, 'david'),
	(7, 'The Client', 2, '978-1-60309-542-6', '', 1, true, 'david'),
	(8, 'The Rainmaker', 2, '978-0-439-02348-1', '', 1, true, 'david'),
	(9, 'Pet Semetary', 1, '978-0-19-853453-6', '', 2, true, 'david'),
	(10, 'Skeleton Crew', 1, '978-3-16-148410-0', '', 2, true, 'david'),
	(11, 'It', 1, '978-1-4028-9462-6', '', 2, true, 'david'),
	(12, 'The Shining', 1, '978-92-95055-02-5', '', 2, true, 'david'),
	(13, 'The Da Vinci Code', 2, '978-1-56619-909-4', '', 1, true, 'david');

INSERT INTO loc.book_author_xref
	(book_author_xref_id, book_id, author_id, created_by)
VALUES
	(1, 1, 2, 'david'),
	(2, 2, 2, 'david'),
	(3, 3, 2, 'david'),
	(4, 4, 1, 'david'),
	(5, 5, 1, 'david'),
	(6, 6, 1, 'david'),
	(7, 7, 1, 'david'),
	(8, 8, 1, 'david'),
	(9, 9, 3, 'david'),
	(10, 10, 3, 'david'),
	(11, 11, 3, 'david'),
	(12, 12, 3, 'david'),
	(13, 13, 4, 'david');

INSERT INTO loc.user
	(user_id, user_name, email, active, created_by)
VALUES
	(1, 'david_galdamez', 'dgaldamez77@gmail.com', true, 'david'),
	(2, 'lynn_keeling', 'LynnKeeling@bluebeam.com', true, 'david'),
	(3, 'layne_cairo', 'LayneCairo@bluebeam.com', true, 'david'),
	(4, 'erika_umana', 'ErikaUmana@bluebeam.com', true, 'david'),
	(5, 'kathy', 'kathy@bluebeam.com', true, 'david');

INSERT INTO loc.book_inventory
	(book_id, book_count, created_by)
SELECT book_id, 10, created 
FROM loc.book;
	
ALTER SEQUENCE loc.book_book_id_seq  START 1;
ALTER SEQUENCE loc.author_author_id_seq START 1;
ALTER SEQUENCE loc.genre_genre_id_seq START 1;
ALTER SEQUENCE loc.book_author_xref_book_author_xref_id_seq START 1;
ALTER SEQUENCE loc.publisher_publisher_id_seq START 1;
ALTER SEQUENCE loc.user_user_id_seq START 1;
