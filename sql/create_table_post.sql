-- Table: public.post

-- DROP TABLE public.post;

CREATE SEQUENCE post_id_seq;
CREATE TABLE public.post
(
  id bigint NOT NULL DEFAULT nextval('post_id_seq'),
  guid character varying(255) COLLATE pg_catalog."default" NOT NULL,
  title character varying(255) COLLATE pg_catalog."default" NOT NULL,
  content text COLLATE pg_catalog."default" NOT NULL,
  created_by bigint NOT NULL,
  created_at timestamp with time zone NOT NULL,
  modified_by bigint,
  modified_at timestamp with time zone,
  CONSTRAINT pk_id PRIMARY KEY (id),
  CONSTRAINT unq_guid UNIQUE (guid)
)
  WITH (
  OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.post
OWNER to flexiportal;
ALTER SEQUENCE post_id_seq OWNED BY post.id;