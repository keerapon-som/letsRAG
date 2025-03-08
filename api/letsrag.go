package api

import "letsrag/entities/db"

type DocumentRepo interface {
	SaveDocument(document string, vector []float64) error
	GetRelatedDocumentsByVector(vector []float64, limit int) ([]db.Document, error)
}

type LetRags struct {
	AI           AIService
	textToVector TextToVector
	repo         DocumentRepo
}

func NewLetsRag(ai AIService, textToVector TextToVector, repo DocumentRepo) *LetRags {
	return &LetRags{
		AI:           ai,
		textToVector: textToVector,
		repo:         repo,
	}
}

func (l *LetRags) SaveDocumentToDB(text string, modelName string) error {
	vector, err := l.textToVector.ConvertTextToVector(text, modelName)
	if err != nil {
		return err
	}

	return l.repo.SaveDocument(text, vector[0])
}

func (l *LetRags) GetRelatedDocuments(text string, modelName string, limit int) ([]string, error) {
	vector, err := l.textToVector.ConvertTextToVector(text, modelName)
	if err != nil {
		return []string{}, err
	}

	docs, err := l.repo.GetRelatedDocumentsByVector(vector[0], limit)
	if err != nil {
		return []string{}, err
	}

	document := []string{}
	for _, doc := range docs {
		document = append(document, doc.Document)
	}

	return document, nil
}
