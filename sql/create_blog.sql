CREATE TABLE blog (
  id BIGINT NOT NULL AUTO_INCREMENT,
  title VARCHAR(500) NOT NULL ,
  content LONGTEXT NOT NULL,
  user_id BIGINT NOT NULL,
  created_at DATETIME NOT NULL,
  modified_by BIGINT NULL,
  modified_at DATETIME NULL,
  PRIMARY KEY (id)
)