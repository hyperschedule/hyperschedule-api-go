CREATE OR REPLACE FUNCTION update_set_timestamp()
RETURNS TRIGGER LANGUAGE 'plpgsql' 
AS $$ BEGIN NEW.updated_at = NOW(); RETURN NEW; END $$;
