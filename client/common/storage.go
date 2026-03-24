package common

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

const (
	MAX_SIZE    = 8192
	AGENCY_PATH = "./agency.csv"
)

type Storage struct {
	file   *os.File
	reader *csv.Reader
	buf    *Bet
}

func NewStorage() (*Storage, error) {
	file, err := os.Open(AGENCY_PATH)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	var buf *Bet = nil

	storage := Storage{file, reader, buf}
	return &storage, nil
}

func (s *Storage) Close() {
	s.file.Close()
}

// Attempts to load `maxAmount` bets from its csv file, if
// the bets exceed `MAX_SIZE` or the end of the file is
// reached, then less bets will be returned.
//
// If `MAX_SIZE` is reached, the exceeding bet is buffered
// so that it can be returned in a later call to the method.
func (s *Storage) LoadBets(maxAmount int) (BetBatch, error) {
	log.Infof("action: load_bets | result: in_progress")
	batch := make([]Bet, 0, maxAmount)
	batch_size := 0
	if s.buf != nil {
		batch = append(batch, *s.buf)
		maxAmount--
		batch_size += s.buf.getSize()
		s.buf = nil
	}

	for i := 0; i < maxAmount; i++ {
		record, err := s.reader.Read()
		if err == io.EOF {
			return batch, err
		}

		if err != nil {
			return nil, err
		}

		bet, record_size, err := s.getBetFromRecord(record)

		if batch_size+record_size > MAX_SIZE {
			s.buf = &bet
			break
		}

		batch = append(batch, bet)
		batch_size += record_size
	}

	return batch, nil
}

func (s *Storage) getBetFromRecord(record []string) (Bet, int, error) {
	bet_fields := 6
	if len(record) != bet_fields {
		return Bet{}, 0, fmt.Errorf(
			"missing fields, there are %v and bet is %v fields",
			len(record), bet_fields)
	}

	Agency := record[0]
	FirstName := record[1]
	LastName := record[2]
	Document := record[3]
	BirthDate := record[4]
	Number := record[5]

	bet := NewBet(Agency, FirstName, LastName, Document, BirthDate, Number)
	betSize := bet.getSize()

	return bet, betSize, nil
}

func (b *Bet) getSize() int {
	betSize := len(b.Agency) + len(b.FirstName) + len(b.LastName) +
		len(b.Document) + len(b.BirthDate) + len(b.Number)

	return betSize
}
