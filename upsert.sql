\connect SQO
DROP FUNCTION merge_cache(integer, varchar(25), varchar(50), jsonb, timestamp);
CREATE OR REPLACE FUNCTION merge_cache(keyId integer, characterId varchar(25), apiPath varchar(50), Data JSONB, cachedUntil timestamp) RETURNS integer AS
$$
DECLARE
    _id integer := 0;
BEGIN
    LOOP
        -- first try to update the key
        EXECUTE 'UPDATE cache SET data=$1 Where keyid = $2 and characterid = $3 and apipath = $4 and cacheduntil = $5 RETURNING id'
         into _id
         USING Data, keyId, characterId, apiPath, cachedUntil;
        IF FOUND or _id != 0 THEN
            RETURN _id;
        END IF;

        -- not there, so try to insert the key
        -- if someone else inserts the same key concurrently,
        -- we could get a unique-key failure
        BEGIN
            EXECUTE 'INSERT INTO cache(keyid, characterid, apipath, data, cacheduntil) VALUES($1, $2, $3, $4, $5) RETURNING id'
             into _id
             USING keyId, characterId, apiPath, Data, cachedUntil;
            RETURN _id;
        EXCEPTION 
            WHEN unique_violation THEN
            -- Do nothing, and loop to try the UPDATE again.
        END;
        
    END LOOP;
END;
$$
LANGUAGE plpgsql;

ALTER DATABASE "SQO" SET Timezone TO 'UTC';