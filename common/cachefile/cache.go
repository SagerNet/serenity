package cachefile

import (
	"errors"
	"os"
	"time"

	"github.com/sagernet/bbolt"
	bboltErrors "github.com/sagernet/bbolt/errors"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
)

var (
	bucketSubscription = []byte("subscription")

	bucketNameList = []string{
		string(bucketSubscription),
	}
)

type CacheFile struct {
	path string
	DB   *bbolt.DB
}

func New(path string) *CacheFile {
	return &CacheFile{
		path: path,
	}
}

func (c *CacheFile) Start() error {
	const fileMode = 0o666
	options := bbolt.Options{Timeout: time.Second}
	var (
		db  *bbolt.DB
		err error
	)
	for i := 0; i < 10; i++ {
		db, err = bbolt.Open(c.path, fileMode, &options)
		if err == nil {
			break
		}
		if errors.Is(err, bboltErrors.ErrTimeout) {
			continue
		}
		if E.IsMulti(err, bboltErrors.ErrInvalid, bboltErrors.ErrChecksum, bboltErrors.ErrVersionMismatch) {
			rmErr := os.Remove(c.path)
			if rmErr != nil {
				return err
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		return err
	}
	err = db.Batch(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bbolt.Bucket) error {
			bucketName := string(name)
			if !(common.Contains(bucketNameList, bucketName)) {
				_ = tx.DeleteBucket(name)
			}
			return nil
		})
	})
	if err != nil {
		db.Close()
		return err
	}
	c.DB = db
	return nil
}

func (c *CacheFile) Close() error {
	if c.DB == nil {
		return nil
	}
	return c.DB.Close()
}

func (c *CacheFile) LoadSubscription(name string) *Subscription {
	var subscription Subscription
	err := c.DB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(bucketSubscription)
		if bucket == nil {
			return nil
		}
		data := bucket.Get([]byte(name))
		if data == nil {
			return nil
		}
		return subscription.UnmarshalBinary(data)
	})
	if err != nil {
		return nil
	}
	return &subscription
}

func (c *CacheFile) StoreSubscription(name string, subscription *Subscription) error {
	data, err := subscription.MarshalBinary()
	if err != nil {
		return err
	}
	return c.DB.Batch(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketSubscription)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(name), data)
	})
}
