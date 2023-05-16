package storage

type Word struct {
	Esp string `json:"esp"`
	Rus string `json:"rus"`
}

type WordRepository interface {
	GetFreeWords() ([]Word, error)
	GetPaidWords() ([]Word, error)
}

type WordRepo struct {
	freeWords []Word
	paidWords []Word
}

func NewWordRepo() *WordRepo {
	freeWords := []Word{
		Word{"chica", "девушка"},
		Word{"padre", "отец"},
		Word{"cuatro", "четыре"},
		Word{"perro", "собака"},
		Word{"hora", "час"},
		Word{"raton", "мышь"},
		Word{"trabajo", "работа"},
	}

	additionalWords := []Word{
		Word{"boligrafo", "ручка"},
		Word{"casa", "дом"},
		Word{"desayuno", "завтрак"},
		Word{"nadar", "плавать"},
		Word{"programador", "программист"},
		Word{"cerveza", "пиво"},
		Word{"bandera en tus manos", "флаг вам в руки"},
	}

	paidWords := append(freeWords, additionalWords...)

	return &WordRepo{
		freeWords: freeWords,
		paidWords: paidWords,
	}
}

func (r *WordRepo) GetFreeWords() ([]Word, error) {
	return r.freeWords, nil
}

func (r *WordRepo) GetPaidWords() ([]Word, error) {
	return r.paidWords, nil
}
