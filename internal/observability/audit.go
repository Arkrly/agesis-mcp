package observability

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

var (
	bucketAudit = []byte("audit")
)

// AuditLogger persists security decisions to BoltDB.
type AuditLogger struct {
	db *bbolt.DB
}

type AuditEntry struct {
	Timestamp  time.Time      `json:"ts"`
	AgentID    string         `json:"agent_id"`
	Method     string         `json:"method"`
	Tool       string         `json:"tool,omitempty"`
	Allowed    bool           `json:"allowed"`
	Reason     string         `json:"reason,omitempty"`
	Inspection map[string]any `json:"inspection,omitempty"`
}

func NewAuditLogger(path string) (*AuditLogger, error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("open audit db: %w", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketAudit)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("create audit bucket: %w", err)
	}

	return &AuditLogger{db: db}, nil
}

func (l *AuditLogger) Log(entry AuditEntry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal audit entry: %w", err)
	}

	return l.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketAudit)
		key := []byte(entry.Timestamp.Format(time.RFC3339Nano))
		return b.Put(key, data)
	})
}

func (l *AuditLogger) List(limit int) ([]AuditEntry, error) {
	var entries []AuditEntry
	err := l.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketAudit)
		c := b.Cursor()

		count := 0
		for k, v := c.Last(); k != nil && count < limit; k, v = c.Prev() {
			var entry AuditEntry
			if err := json.Unmarshal(v, &entry); err != nil {
				continue
			}
			entries = append(entries, entry)
			count++
		}
		return nil
	})
	return entries, err
}

func (l *AuditLogger) Close() error {
	return l.db.Close()
}
