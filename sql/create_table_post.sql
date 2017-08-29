-- Table: public.post

-- DROP TABLE public.post;

CREATE TABLE public.post
(
  id BIGSERIAL,
  guid character varying(255) NOT NULL UNIQUE,
  title character varying(255) NOT NULL,
  content text COLLATE pg_catalog."default" NOT NULL,
  created_by bigint NOT NULL,
  created_at timestamp with time zone NOT NULL,
  modified_by bigint,
  modified_at timestamp with time zone,
  CONSTRAINT pk_id PRIMARY KEY (id),
  CONSTRAINT unq_guid UNIQUE (guid)
)