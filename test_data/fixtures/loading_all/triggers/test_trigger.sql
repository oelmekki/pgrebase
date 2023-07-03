CREATE FUNCTION test_trigger2()
RETURNS trigger
LANGUAGE plpgsql
AS $$
  BEGIN
    NEW.active := true;
    RETURN NEW;
  END
$$;

CREATE TRIGGER test_trigger2 BEFORE INSERT ON users
FOR EACH ROW EXECUTE PROCEDURE test_trigger2();
