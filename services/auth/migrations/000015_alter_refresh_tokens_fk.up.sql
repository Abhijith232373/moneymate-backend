DO $$ 
DECLARE 
    fk_name text;
BEGIN
    SELECT conname INTO fk_name
    FROM pg_constraint
    WHERE conrelid = 'auth.refresh_tokens'::regclass
      AND contype = 'f';

    IF fk_name IS NOT NULL THEN
        EXECUTE 'ALTER TABLE auth.refresh_tokens DROP CONSTRAINT ' || fk_name;
    END IF;
END $$;
