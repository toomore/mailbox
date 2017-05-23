CREATE TABLE IF NOT EXISTS campaign (
    id CHAR(8) NOT NULL UNIQUE,
    seed CHAR(8) NOT NULL,
    created DATETIME DEFAULT NOW(),
    updated DATETIME DEFAULT NOW() ON UPDATE NOW()
    ) CHARACTER SET 'utf8' 
      COLLATE 'utf8_icelandic_ci';

CREATE TABLE IF NOT EXISTS user (
    id INT NOT NULL AUTO_INCREMENT,
    email VARCHAR(64) NOT NULL,
    groups VARCHAR(64) NOT NULL,
    f_name VARCHAR(32) NOT NULL,
    l_name VARCHAR(32) not NULL,
    created DATETIME DEFAULT NOW(),
    updated DATETIME DEFAULT NOW() ON UPDATE NOW(),
    alive BOOL DEFAULT true,
    PRIMARY KEY (id),
    UNIQUE KEY email_groups (email, groups),
    INDEX groups (groups ASC)
    ) CHARACTER SET 'utf8'
      COLLATE 'utf8_icelandic_ci';

CREATE TABLE IF NOT EXISTS reader (
    no INT NOT NULL AUTO_INCREMENT,
    uid INT NOT NULL,
    cid CHAR(8) NOT NULL,
    ip CHAR(15) NOT NULL,
    agent CHAR(128) NOT NULL,
    created DATETIME DEFAULT NOW(),
    PRIMARY KEY (no),
    INDEX uid (uid ASC),
    INDEX cid (cid ASC)
    ) CHARACTER SET 'utf8'
      COLLATE 'utf8_icelandic_ci';

CREATE TABLE IF NOT EXISTS links (
    id CHAR(8) NOT NULL,
    cid CHAR(8) NOT NULL,
    url CHAR(255) NOT NULL,
    urlhash CHAR(32) NOT NULL,
    created DATETIME DEFAULT NOW(),
    UNIQUE KEY cid_urlhash (cid, urlhash),
    INDEX cid (cid ASC)
    ) CHARACTER SET 'utf8'
      COLLATE 'utf8_icelandic_ci';
