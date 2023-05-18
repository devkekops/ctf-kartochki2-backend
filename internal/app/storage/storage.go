package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Word struct {
	Esp string `json:"esp" db:"esp"`
	Rus string `json:"rus" db:"rus"`
}

type WordRepository interface {
	GetFreeWords() ([]Word, error)
	GetPaidWords() ([]Word, error)
}

type WordRepo struct {
	freeWords []Word
	paidWords []Word
}

func NewWordRepo(dsn string) (*WordRepo, error) {
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	freeWords := []Word{}
	err = db.Select(&freeWords, `SELECT * FROM freeWords`)
	if err != nil {
		return nil, err
	}

	paidWords := []Word{}
	err = db.Select(&paidWords, `SELECT * FROM paidWords`)
	if err != nil {
		return nil, err
	}

	return &WordRepo{
		freeWords: freeWords,
		paidWords: paidWords,
	}, nil
}

func (r *WordRepo) GetFreeWords() ([]Word, error) {
	return r.freeWords, nil
}

func (r *WordRepo) GetPaidWords() ([]Word, error) {
	return r.paidWords, nil
}
