CREATE TABLE articles (
	id     	SERIAL,
	author_id	integer,
    title    	varchar(40),
	body		text,
    PRIMARY KEY(id)
);
