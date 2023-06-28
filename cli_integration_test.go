package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

var dbConnectionScheme string

func TestMain(m *testing.M) {

	start := exec.Command("./test_data/reset_db.sh")
	err := start.Run()
	if err != nil {
		fmt.Println("Can't start the database")
		os.Exit(1)
	}

	exitVal := 1

	defer (func() {
		stop := exec.Command("./test_data/stop_db.sh")
		out, err := stop.Output()
		if err != nil {
			fmt.Printf("%s\n", out)
			os.Exit(1)
		}

		os.Exit(exitVal)
	})()

	port := os.Getenv("PG_TEST_PORT")
	if len(port) == 0 {
		port = "5433"
	}
	dbConnectionScheme = fmt.Sprintf("postgres://postgres:@localhost:%s/pgrebase?sslmode=disable", port)
	os.Setenv("DATABASE_URL", dbConnectionScheme)

	// We want to test the binary directly, so let's build it.
	// for testing the functions, see core package.
	start = exec.Command("go", "build")
	out, err := start.Output()
	if err != nil {
		fmt.Printf("%s\n", out)
	} else {
		exitVal = m.Run()
	}
}

func query(query string, parameters ...interface{}) (Rows *sql.Rows, err error) {
	var co *sql.DB
	co, err = sql.Open("postgres", dbConnectionScheme)
	if err != nil {
		err = fmt.Errorf("can't connect to db : %v", err)
		return
	}
	defer co.Close()

	Rows, err = co.Query(query, parameters...)
	if err != nil {
		err = fmt.Errorf("error while executing query : %v", err)
		return
	}

	return
}

func TestLoadingAFunction(t *testing.T) {
	start := exec.Command("./pgrebase", "test_data/fixtures/loading_a_function/")
	out, err := start.Output()
	if err != nil {
		fmt.Printf("%s\n", out)
		t.Fatalf("Error while executing pgrebase: %v", err)
	}

	rows, err := query("SELECT test_function()")
	if err != nil {
		t.Fatalf("Can't query : %v", err)
	}

	defer rows.Close()

	if !rows.Next() {
		t.Fatalf("Calling function does not provide any result.")
		return
	}
}

func TestLoadingAView(t *testing.T) {
	t.Cleanup(func() {
		rows, err := query("DELETE FROM users")
		if err != nil {
			t.Fatalf("Can't query : %v", err)
		}
		rows.Close()
	})

	start := exec.Command("./pgrebase", "test_data/fixtures/loading_a_view/")
	out, err := start.Output()
	if err != nil {
		fmt.Printf("%s\n", out)
		t.Fatalf("Error while executing pgrebase: %v", err)
	}

	rows, err := query("INSERT INTO users(name, bio) VALUES('John Doe', 'John Doe does stuff.')")
	if err != nil {
		fmt.Printf("Can't create mock record : %v\n", err)
		t.Fatalf("Can't insert test record.")
	}
	rows.Close()

	rows, err = query("SELECT * FROM test_view")
	if err != nil {
		t.Fatalf("Can't query : %v", err)
	}

	defer rows.Close()

	if !rows.Next() {
		t.Fatalf("Calling function does not provide any result.")
		return
	}

	id := 0
	name := ""
	err = rows.Scan(&id, &name)
	if err != nil {
		t.Fatalf("Can't fetch columns : %v", err)
	}

	if id != 1 {
		t.Errorf("ID 1 expected, got %d", id)
	}

	if name != "John Doe" {
		t.Errorf("Name \"John Doe\" expected, got %s", name)
	}
}

func TestLoadingATrigger(t *testing.T) {
	t.Cleanup(func() {
		rows, err := query("DELETE FROM users")
		if err != nil {
			t.Fatalf("Can't query : %v", err)
		}
		rows.Close()
	})

	start := exec.Command("./pgrebase", "test_data/fixtures/loading_a_trigger/")
	out, err := start.Output()
	if err != nil {
		fmt.Printf("%s\n", out)
		t.Fatalf("Error while executing pgrebase: %v", err)
	}

	rows, err := query("INSERT INTO users(name, bio) VALUES('John Doe', 'John Doe does stuff.')")
	if err != nil {
		fmt.Printf("Can't create mock record : %v\n", err)
		t.Fatalf("Can't insert test record.")
	}
	rows.Close()

	rows, err = query("SELECT active FROM users")
	if err != nil {
		t.Fatalf("Can't query : %v", err)
	}

	defer rows.Close()

	if !rows.Next() {
		t.Fatalf("Calling function does not provide any result.")
		return
	}

	active := false
	err = rows.Scan(&active)
	if err != nil {
		t.Fatalf("Can't fetch columns : %v", err)
	}

	if !active {
		t.Errorf("Trigger expected to set `active` to true, it's false.")
	}
}

func TestLoadingAllTypes(t *testing.T) {
	start := exec.Command("./pgrebase", "test_data/fixtures/loading_all/")
	out, err := start.Output()
	if err != nil {
		fmt.Printf("%s\n", out)
		t.Fatalf("Error while executing pgrebase: %v", err)
	}
}

func TestLoadingWithDependencies(t *testing.T) {
	start := exec.Command("./pgrebase", "test_data/fixtures/dependencies/")
	out, err := start.Output()
	if err != nil {
		fmt.Printf("%s\n", out)
		t.Fatalf("Error while executing pgrebase: %v", err)
	}
}

func TestLoadingWithWatcher(t *testing.T) {
	go (func() {
		start := exec.Command("./pgrebase", "test_data/fixtures/watcher/")
		out, err := start.Output()
		if err != nil {
			fmt.Printf("%s\n", out)
			fmt.Printf("Error while executing pgrebase: %v", err)
			os.Exit(1)
		}
	})()

	time.Sleep(1 * time.Second)

	rows, err := query("INSERT INTO users(name, bio) VALUES('John Doe', 'John Doe does stuff.')")
	if err != nil {
		t.Fatalf("Can't insert test record : %v.", err)
	}
	rows.Close()

	testFile := "./test_data/fixtures/watcher/views/test_view5.sql"
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Can't create test view file : %v.", err)
	}

	_, err = fmt.Fprintf(file, "CREATE VIEW test_view5 AS SELECT id, name FROM users;")
	if err != nil {
		t.Fatalf("Can't write test view in file : %v.", err)
	}

	file.Close()
	t.Cleanup(func() {
		os.Remove(testFile)
	})

	time.Sleep(1 * time.Second)

	rows, err = query("INSERT INTO users(name, bio) VALUES('John Doe', 'John Doe does stuff.')")
	if err != nil {
		fmt.Printf("Can't create mock record : %v\n", err)
		t.Fatalf("Can't insert test record.")
	}
	rows.Close()

	rows, err = query("SELECT * FROM test_view")
	if err != nil {
		t.Fatalf("Can't query : %v", err)
	}

	defer rows.Close()

	if !rows.Next() {
		t.Fatalf("Calling function does not provide any result.")
		return
	}

	id := 0
	name := ""
	err = rows.Scan(&id, &name)
	if err != nil {
		t.Fatalf("Can't fetch columns : %v", err)
	}

	if name != "John Doe" {
		t.Errorf("Name \"John Doe\" expected, got %s", name)
	}
}
