package db

import (
	"encoding/binary"
	"time"

	"go.etcd.io/bbolt"
)

var taskBucket = []byte("tasks")
var db *bbolt.DB

type Task struct {
	Key   int
	Value string
}

func Init(dbPath string) error {
	var err error
	db, err = bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func CreateTask(task string) (int, error) {
	var id int
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task))

	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(taskBucket)
		b.ForEach(func(k, v []byte) error {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteTask(id int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(id))
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))

}
