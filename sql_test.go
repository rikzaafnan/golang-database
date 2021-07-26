package golangdatabase

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {

	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(id, name) VALUES ('joko', 'Joko')"

	_, err := db.ExecContext(ctx, script)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")

}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name FROM customer"

	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	// fmt.Println("Success insert new customer")

	for rows.Next() {

		var id, name string

		err = rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println("Id:", id)
		fmt.Println("Name:", name)

	}

}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at  FROM customer"

	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {

		var id, name string
		var email sql.NullString
		var balance int32
		var rating float32
		var birthDate sql.NullTime
		var createdAt time.Time
		var married bool

		err = rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("=====================")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)

		if email.Valid {
			fmt.Println("Email:", email)

		}

		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)

		if birthDate.Valid {
			fmt.Println("Birth Date:", birthDate)
		}

		fmt.Println("Married:", married)
		fmt.Println("Created At:", createdAt)
		fmt.Println("=====================")

	}
}
func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin'; #"
	password := "salah"

	script := "SELECT username FROM user where username = '" + username + "' AND password='" + password + "' Limit 1 "

	fmt.Println(script)

	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {

		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Sukses Login", username)

	} else {
		fmt.Println("Gagal Login")
	}

}

func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	// username := "admin'; #"
	username := "admin"
	password := "admin"

	script := "SELECT username FROM user where username = ? AND password = ? Limit 1 "

	// fmt.Println(script)

	rows, err := db.QueryContext(ctx, script, username, password)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {

		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Sukses Login", username)

	} else {
		fmt.Println("Gagal Login")
	}

}

func TestExecSqlParameter(t *testing.T) {

	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "eko"
	password := "eko"

	script := "INSERT INTO user(username, password) VALUES (?, ?)"

	_, err := db.ExecContext(ctx, script, username, password)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new user")

}

func TestAutoIncrement(t *testing.T) {

	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "eko@gmail.com"
	comment := "Test komen"

	script := "INSERT INTO comments(email, comment) VALUES (?, ?)"

	result, err := db.ExecContext(ctx, script, email, comment)

	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert comment with id", insertId)

}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO comments(email, comment) VALUES (?, ?)"

	statement, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	for i := 0; i < 10; i++ {

		email := "eko" + strconv.Itoa(i) + "@gmail.com"
		comment := "Komentar ke " + strconv.Itoa(i)

		result, err := statement.ExecContext(ctx, email, comment)

		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id", id)

	}

}

func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	script := "INSERT INTO comments(email, comment) VALUES (?, ?)"

	// db transaction
	for i := 0; i < 10; i++ {

		email := "eko" + strconv.Itoa(i) + "@gmail.com"
		comment := "Komentar ke " + strconv.Itoa(i)

		result, err := tx.ExecContext(ctx, script, email, comment)

		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id", id)

	}

	// err = tx.Commit()
	err = tx.Rollback()

	if err != nil {
		panic(err)
	}
}
