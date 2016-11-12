package main

/*
 * Usage :
 *   DATABASE_URL=url pgrebase [-w] sql_dir/
 */
func ParseConfig() {
}

/*
 * Expected target structure:
 *
 * <sql_dir>/
 *   functions/
 *   triggers/
 *   views/
 *
 * At least one of functions/triggers/views/ should exist.
 *
 */
func CheckSanity() {
}

/*
 * Start the actual work
 */
func Process() {
}

func main() {
	ParseConfig()
	CheckSanity()
	Process()
}
