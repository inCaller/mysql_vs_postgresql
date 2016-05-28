DROP DATABASE IF EXISTS mysql_vs_pgsql;
CREATE DATABASE mysql_vs_pgsql;

CREATE TABLE mysql_vs_pgsql.users (
	user_id   BIGINT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
	user_name VARCHAR(64) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE UNIQUE INDEX users_user_name ON mysql_vs_pgsql.users (user_name);

CREATE TABLE mysql_vs_pgsql.friends (
	user_id   BIGINT,
	friend_id BIGINT,
	FOREIGN KEY (user_id)   REFERENCES mysql_vs_pgsql.users (user_id),
	FOREIGN KEY (friend_id) REFERENCES mysql_vs_pgsql.users (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE INDEX friends_user_id   ON mysql_vs_pgsql.friends (user_id);
-- CREATE INDEX friends_friend_id ON mysql_vs_pgsql.friends (friend_id);

CREATE TABLE mysql_vs_pgsql.messages (
	msg_id    BIGINT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
	user_id   BIGINT,
	ctime     TIMESTAMP NOT NULL,
	message   VARCHAR(16384) NOT NULL,
	FOREIGN KEY (user_id)   REFERENCES mysql_vs_pgsql.users (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE INDEX friends_user_id   ON mysql_vs_pgsql.messages (user_id);
