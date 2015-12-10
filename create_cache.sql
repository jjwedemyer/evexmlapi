-- Table: cache

-- DROP TABLE cache;

CREATE TABLE cache
(
  keyid integer,
  apipath character varying(50),
  id integer NOT NULL DEFAULT nextval('skillqueue_id_seq'::regclass),
  characterid character varying(25),
  vcode character varying(100),
  data jsonb,
  cacheduntil timestamp with time zone NOT NULL,
  CONSTRAINT id PRIMARY KEY (id),
  CONSTRAINT key UNIQUE (keyid, apipath, characterid, cacheduntil)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE cache
  OWNER TO postgres;
