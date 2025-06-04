package history

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

type Entry struct {
	ID        uint64
	Timestamp time.Time
	Command   string
	Query     string
	Response  string
	Selected  *string
}

type History struct {
	db *bbolt.DB
}

const (
	bucketName = "history"
	maxEntries = 5
)

func GetHistoryPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".ted", "history.db"), nil
}

func Load() (*History, error) {
	dbPath, err := GetHistoryPath()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create history directory: %w", err)
	}

	db, err := bbolt.Open(dbPath, 0644, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open history database: %w", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create history bucket: %w", err)
	}

	return &History{db: db}, nil
}

func (h *History) Close() error {
	if h.db != nil {
		return h.db.Close()
	}
	return nil
}

func (h *History) AddEntry(command, query, response string, selected *string) error {
	return h.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		id, err := bucket.NextSequence()
		if err != nil {
			return fmt.Errorf("failed to generate entry ID: %w", err)
		}

		entry := Entry{
			ID:        id,
			Timestamp: time.Now(),
			Command:   command,
			Query:     query,
			Response:  response,
			Selected:  selected,
		}

		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		if err := encoder.Encode(entry); err != nil {
			return fmt.Errorf("failed to encode entry: %w", err)
		}
		data := buf.Bytes()
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, id)

		if err := bucket.Put(key, data); err != nil {
			return fmt.Errorf("failed to store entry: %w", err)
		}

		return h.trimToMaxEntries(bucket)
	})
}

func (h *History) GetEntries() ([]Entry, error) {
	var entries []Entry

	err := h.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return nil
		}

		cursor := bucket.Cursor()

		for k, v := cursor.Last(); k != nil; k, v = cursor.Prev() {
			var entry Entry
			buf := bytes.NewBuffer(v)
			decoder := gob.NewDecoder(buf)
			if err := decoder.Decode(&entry); err != nil {
				continue
			}
			entries = append(entries, entry)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve entries: %w", err)
	}

	return entries, nil
}

func (h *History) DeleteMostRecent() error {
	return h.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("no entries to delete")
		}

		cursor := bucket.Cursor()
		k, _ := cursor.Last()
		if k == nil {
			return fmt.Errorf("no entries to delete")
		}

		return bucket.Delete(k)
	})
}

func (h *History) Clear() error {
	return h.db.Update(func(tx *bbolt.Tx) error {
		if err := tx.DeleteBucket([]byte(bucketName)); err != nil {
			return fmt.Errorf("failed to delete history bucket: %w", err)
		}

		_, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("failed to recreate history bucket: %w", err)
		}

		return nil
	})
}

func (h *History) trimToMaxEntries(bucket *bbolt.Bucket) error {
	cursor := bucket.Cursor()

	count := 0
	for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
		count++
	}

	if count <= maxEntries {
		return nil
	}

	toDelete := count - maxEntries
	cursor = bucket.Cursor()

	for k, _ := cursor.First(); k != nil && toDelete > 0; k, _ = cursor.Next() {
		if err := bucket.Delete(k); err != nil {
			return fmt.Errorf("failed to delete old entry: %w", err)
		}
		toDelete--
	}

	return nil
}
