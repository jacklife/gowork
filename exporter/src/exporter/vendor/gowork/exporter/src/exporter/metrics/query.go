package metrics

import (
	"database/sql"
	"log"
)

func QueryVersion(db *sql.DB)string  {
	rows, err := db.Query("show server_version")
	if err != nil {
		log.Fatal(err)
	}
	var version string
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatal(err)
		}
	}

	return version

}

func queryTableCount(db *sql.DB) (float64, error) {
	rows, err := db.Query("select count(*) from pg_tables;")
	if err != nil {
		return 0, err
	}

	var tableCount float64
	for rows.Next() {
		err := rows.Scan(&tableCount)
		if err != nil {
			return 0, err
		}
	}

	return tableCount, nil
}

func queryTotalSize(db *sql.DB) (float64, error) {
	rows, err := db.Query("select sum(pg_database_size(pg_database.datname)) as size from pg_database;")
	if err != nil {
		return 0, err
	}

	var totalSize float64
	for rows.Next() {
		err := rows.Scan(&totalSize)
		if err != nil {
			return 0, err
		}
	}

	return totalSize, nil
}

func queryMaxConnections(db *sql.DB) (float64, error) {
	rows, err := db.Query("show max_connections;")
	if err != nil {
		return 0, err
	}

	var maxConn float64
	for rows.Next() {
		err := rows.Scan(&maxConn)
		if err != nil {
			return 0, err
		}
	}

	return maxConn, nil
}

func queryCurrentConnections(db *sql.DB) (float64, error) {
	rows, err := db.Query("select count(*) from pg_stat_activity;")
	if err != nil {
		return 0, err
	}

	var currentConn float64
	for rows.Next() {
		err := rows.Scan(&currentConn)
		if err != nil {
			return 0, err
		}
	}

	return currentConn, nil
}

func queryActiveConnections(db *sql.DB) (float64, error) {
	rows, err := db.Query("select count(*)  from pg_stat_activity where not pid=pg_backend_pid();")
	if err != nil {
		return 0, err
	}

	var activeConn float64
	for rows.Next() {
		err := rows.Scan(&activeConn)
		if err != nil {
			return 0, err
		}
	}

	return activeConn, nil
}

func queryDBSize(db *sql.DB) (map[string]float64, error) {
	rows, err := db.Query("select pg_database.datname, pg_database_size(pg_database.datname) as size from pg_database;")
	if err != nil {
		return nil, err
	}
	var dbs = make(map[string]float64)
	for rows.Next() {
		var dbName string
		var dbSize float64
		err := rows.Scan(&dbName, &dbSize)

		if err != nil {
			return nil, err
		}

		dbs[dbName] = dbSize
	}

	return dbs, nil
}
