CREATE TABLE IF NOT EXISTS campaign(
  id char(8) NOT NULL UNIQUE,
  seed char(8) NOT NULL,
  created DATETIME DEFAULT NOW(),
  updated DATETIME DEFAULT NOW() ON UPDATE NOW()) character
SET 'utf8' COLLATE 'utf8_icelandic_ci';

CREATE TABLE IF NOT EXISTS user (
  id int NOT NULL AUTO_INCREMENT,
  email varchar(64) NOT NULL,
  email_uni varchar(64) NOT NULL,
  groups VARCHAR(64) NOT NULL,
  f_name varchar(32) NOT NULL,
  l_name varchar(32) NOT NULL,
  created DATETIME DEFAULT NOW(),
  updated DATETIME DEFAULT NOW() ON UPDATE NOW(),
  alive bool DEFAULT TRUE,
  PRIMARY KEY (id),
  UNIQUE KEY email_groups(email_uni, GROUPS),
  INDEX GROUPS (GROUPS ASC)) character
SET 'utf8' COLLATE 'utf8_icelandic_ci';

CREATE TABLE IF NOT EXISTS reader(
  no INT NOT NULL AUTO_INCREMENT,
  uid int NOT NULL,
  cid char(8) NOT NULL,
  ip char(15) NOT NULL,
  agent char(255) NOT NULL,
  created DATETIME DEFAULT NOW(),
  PRIMARY KEY (NO),
  INDEX uid(uid ASC),
  INDEX cid(cid ASC)) character
SET 'utf8' COLLATE 'utf8_icelandic_ci';

CREATE TABLE IF NOT EXISTS doors(
  no INT NOT NULL AUTO_INCREMENT,
  uid int NOT NULL,
  cid char(8) NOT NULL,
  linkid char(8) NOT NULL,
  ip char(15) NOT NULL,
  agent char(255) NOT NULL,
  created DATETIME DEFAULT NOW(),
  PRIMARY KEY (NO),
  INDEX uid(uid ASC),
  INDEX linkid(linkid ASC),
  INDEX cid(cid ASC)) character
SET 'utf8' COLLATE 'utf8_icelandic_ci';

CREATE TABLE IF NOT EXISTS links(
  id char(8) NOT NULL,
  cid char(8) NOT NULL,
  url text NOT NULL,
  urlhash char(32) NOT NULL,
  created DATETIME DEFAULT NOW(),
  UNIQUE KEY cid_urlhash(cid, urlhash),
  INDEX cid(cid ASC)) character
SET 'utf8' COLLATE 'utf8_icelandic_ci';

CREATE TABLE IF NOT EXISTS vote(
  no INT NOT NULL AUTO_INCREMENT,
  id char(8) NOT NULL,
  ip char(15) NOT NULL,
  agent char(255) NOT NULL,
  created DATETIME DEFAULT NOW(),
  PRIMARY KEY (NO),
  INDEX id(id ASC)) character
SET 'utf8' COLLATE 'utf8_icelandic_ci';

