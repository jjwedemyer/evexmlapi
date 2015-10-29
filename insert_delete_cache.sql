\connect SQO
DROP FUNCTION insert_delete_cache(integer, varchar(25), varchar(50), jsonb, timestamp);
CREATE OR REPLACE FUNCTION insert_delete_cache(keyId integer, characterId varchar(25), apiPath varchar(50), Data JSONB, cachedUntil timestamp) RETURNS integer AS
$$
DECLARE
    _id integer := 0;
BEGIN
    
EXECUTE 'INSERT INTO cache(keyid, characterid, apipath, data, cacheduntil) VALUES($1, $2, $3, $4, $5) RETURNING id'
     into _id
     USING keyId, characterId, apiPath, Data, cachedUntil;
EXECUTE 'DELETE From cache Where keyid = $1 and characterid = $2 and apipath = $3 and cachedUntil < $4'
  USING keyId, characterId, apiPath, cachedUntil;
    RETURN _id;
EXCEPTION 
    WHEN unique_violation THEN
        RETURN _id;
    
END;
$$
LANGUAGE plpgsql;