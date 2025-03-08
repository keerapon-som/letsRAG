package repository

import (
	"database/sql"
	"fmt"
	"letsrag/entities/db"
	"letsrag/postgresql"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

type documentRepoPosgresql struct {
	pgDB *sql.DB
}

func NewDocumentRepoPostgresql() *documentRepoPosgresql {
	return &documentRepoPosgresql{
		pgDB: postgresql.DB(),
	}
}

func (d *documentRepoPosgresql) CreateDocumentTable() error { // VECTOR 384 ขึ้นอยู่กับว่า ตอน Gen Text To Vector กำหนด Dimension ไว้เท่าไหร่
	query := `
	CREATE EXTENSION IF NOT EXISTS vector;
	CREATE TABLE IF NOT EXISTS documents (
		id SERIAL PRIMARY KEY,
		document TEXT,
		vector VECTOR(384) -- Adjust the dimension as needed
	)`
	_, err := d.pgDB.Exec(query)
	return err
}

func (d *documentRepoPosgresql) SaveDocument(document string, vector []float64) error {
	// Convert the vector to a string in the format expected by the VECTOR type
	vectorStr := fmt.Sprintf("[%s]", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(vector)), ","), "[]"))
	query := `
    INSERT INTO documents (document, vector)
    VALUES ($1, $2::vector)`
	_, err := d.pgDB.Exec(query, document, vectorStr)
	return err
}

func (d *documentRepoPosgresql) GetRelatedDocumentsByVector(vector []float64, limit int) ([]db.Document, error) {
	// Convert the vector to a string in the format expected by the VECTOR type
	vectorStr := fmt.Sprintf("[%s]", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(vector)), ","), "[]"))
	query := `
    SELECT document, vector
    FROM documents
    ORDER BY vector <-> $1::vector
    LIMIT $2`
	rows, err := d.pgDB.Query(query, vectorStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []db.Document
	for rows.Next() {
		var doc db.Document
		var vectorStr string
		if err := rows.Scan(&doc.Document, &vectorStr); err != nil {
			return nil, err
		}
		// Convert the vector string back to a slice of float64
		vectorStr = strings.Trim(vectorStr, "[]")
		vectorFields := strings.Split(vectorStr, ",")
		vector := make([]float64, len(vectorFields))
		for i, v := range vectorFields {
			vector[i], err = strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
		}
		doc.Vector = vector
		documents = append(documents, doc)
	}
	return documents, nil
}

func (d *documentRepoPosgresql) DeleteDocuments(index []int) error {
	query := `
	DELETE FROM documents
	WHERE id = ANY($1)`
	_, err := d.pgDB.Exec(query, pq.Array(index))
	return err
}

func (d *documentRepoPosgresql) EditDocuments(docs []db.Document) error {
	for _, doc := range docs {
		query := `
		UPDATE documents
		SET document = $1, vector = $2
		WHERE id = $3`
		_, err := d.pgDB.Exec(query, doc.Document, pq.Array(doc.Vector), doc.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *documentRepoPosgresql) GetDocumentsByIndex(index []int) ([]db.Document, error) {
	query := `
	SELECT document, vector
	FROM documents
	WHERE id = ANY($1)`
	rows, err := d.pgDB.Query(query, pq.Array(index))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []db.Document
	for rows.Next() {
		var doc db.Document
		if err := rows.Scan(&doc.Document, &doc.Vector); err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func (d *documentRepoPosgresql) GetDocumentsByRangeOfIndex(start, end int) ([]db.Document, error) {
	query := `
	SELECT document, vector
	FROM documents
	WHERE id BETWEEN $1 AND $2`
	rows, err := d.pgDB.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []db.Document
	for rows.Next() {
		var doc db.Document
		if err := rows.Scan(&doc.Document, &doc.Vector); err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

// Assuming you have a function to convert text to vector
func textToVector(text string) []float64 {
	// Implement your text to vector conversion logic here
	return []float64{0.0, 0.0, 0.0} // Example vector
}
