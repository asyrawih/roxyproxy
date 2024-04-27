package parser

import (
	"fmt"
	"strings"
)

type ConnectionDetails struct {
	User            string
	Database        string
	ApplicationName string
	ClientEncoding  string
}

type SQLQuery struct {
	QueryType string
	Columns   []string
}

func ParseData(data string) {
	var sqlQuery SQLQuery

	// Split the query string by whitespace
	parts := strings.Fields(data)

	// The first part is the query type (e.g., SELECT, INSERT, UPDATE, etc.)
	if len(parts) > 0 {
		sqlQuery.QueryType = strings.ToUpper(parts[0])
	}

	// The remaining parts are the columns or other clauses
	if len(parts) > 1 {
		sqlQuery.Columns = parts[1:]
	}

	s := rebuildSQLQuery(sqlQuery)
	fmt.Printf("query %v\n", s)

	// Parse SQL query
}

func rebuildSQLQuery(sqlQuery SQLQuery) string {
	// Reconstruct the SQL query based on the parsed information
	query := sqlQuery.QueryType + " " + strings.Join(sqlQuery.Columns, " ")
	query = strings.ReplaceAll(query, "Q", "")
	query = strings.ReplaceAll(query, ";", "")
	query = strings.ReplaceAll(query, "#", "")
	query = strings.ReplaceAll(query, "$", "")
	query = strings.ReplaceAll(query, "%", "")
	query = strings.ReplaceAll(query, "+", "")
	return query + ";"
}
