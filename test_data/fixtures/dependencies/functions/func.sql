-- require "test_function3_1.sql"
CREATE FUNCTION test_function3()
RETURNS int
LANGUAGE plpgsql
AS $$
  DECLARE
    
  BEGIN
    SELECT test_function3_1();
    RETURN 1;
  END
$$

