package storage

type Word struct {
	Esp string `json:"esp"`
	Rus string `json:"rus"`
}

type WordRepository interface {
	GetWords() ([]Word, error)
}

type WordRepo struct {
	words []Word
}

func NewWordRepo() *WordRepo {
	var words = []Word{
		Word{"chica", "девушка"},
		Word{"padre", "отец"},
		Word{"cuatro", "четыре"},
		Word{"perro", "собака"},
		Word{"hora", "час"},
		Word{"raton", "мышь"},
		Word{"trabajo", "работа"},
	}
	return &WordRepo{
		words: words,
	}
}

func (r *WordRepo) GetWords() ([]Word, error) {
	return r.words, nil
}
