CREATE TABLE users (
	user_id   SERIAL  NOT NULL PRIMARY KEY,
	user_name VARCHAR(64) NOT NULL
);

CREATE UNIQUE INDEX users_user_name ON users (user_name);

CREATE TABLE friends (
	user_id   BIGINT,
	friend_id BIGINT,
	FOREIGN KEY (user_id)   REFERENCES users (user_id),
	FOREIGN KEY (friend_id) REFERENCES users (user_id)
);

CREATE INDEX friends_user_id   ON friends (user_id);
CREATE UNIQUE INDEX friends_user_id_friend_id ON friends (user_id, friend_id);
-- CREATE INDEX friends_friend_id ON mysql_vs_pgsql.friends (friend_id);

CREATE TABLE messages (
	msg_id    SERIAL  NOT NULL PRIMARY KEY,
	user_id   BIGINT,
	ctime     TIMESTAMP NOT NULL,
	message   VARCHAR(16384) NOT NULL,
	FOREIGN KEY (user_id)   REFERENCES users (user_id)
);

CREATE INDEX messages_user_id   ON messages (user_id);
CREATE INDEX messages_ctime     ON messages (ctime);
