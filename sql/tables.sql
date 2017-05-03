CREATE TABLE IF NOT EXISTS campaign (
    id CHAR(8) NOT NULL UNIQUE,
    seed CHAR(8) NOT NULL,
    created DATETIME DEFAULT NOW(),
    updated DATETIME DEFAULT NOW() ON UPDATE NOW()
    ) CHARACTER SET 'utf8' 
      COLLATE 'utf8_icelandic_ci';

