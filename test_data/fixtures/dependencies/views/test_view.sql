-- require "test_view3_1.sql"
CREATE VIEW test_view3 AS
SELECT id, name FROM users UNION SELECT * FROM test_view3_1;
