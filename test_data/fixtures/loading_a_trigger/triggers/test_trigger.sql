CREATE FUNCTION test_trigger()
RETURNS trigger
LANGUAGE plpgsql
AS $$
  BEGIN
    NEW.active := true;
    RETURN NEW;
  END
$$;

CREATE TRIGGER test_trigger BEFORE INSERT ON users
FOR EACH ROW EXECUTE PROCEDURE test_trigger();
