CREATE TABLE IF NOT EXISTS books
(
	id   varchar(25) PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	CHECK ( length(name) <= 50 )
);

CREATE TABLE IF NOT EXISTS transaction_types
(
	id               INTEGER PRIMARY KEY,
	name          text    NOT NULL,
	CHECK ( length(name) <= 50 )
);

INSERT INTO transaction_types (id, name)
VALUES (1, 'expenses'), (2, 'income')
;

CREATE TABLE IF NOT EXISTS transactions
(
	id               varchar(25) PRIMARY KEY,
	book_id          varchar(25)    NOT NULL REFERENCES "books" (id),
	description      TEXT           NOT NULL,
	amount           numeric(15, 2) NOT NULL DEFAULT 0.0,
	type             INTEGER            NOT NULL REFERENCES "transaction_types" (id) DEFAULT 1,
	date_transaction datetime not null,
	date_created datetime not null,

	CHECK ( length(description) <= 300 )
);