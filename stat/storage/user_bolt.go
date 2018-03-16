package storage

import (
	"log"
	"strings"

	"github.com/boltdb/bolt"
)

const (
	KYC_CATEGORY string = "0x0000000000000000000000000000000000000000000000000000000000000004"

	CATLOG_PROCESSOR_STATE string = "catlog_processor_state"

	ADDRESS_CATEGORY  string = "address_category"
	ADDRESS_ID        string = "address_id"
	ID_ADDRESSES      string = "id_addresses"
	ADDRESS_TIME      string = "address_time"
	PENDING_ADDRESSES string = "pending_addresses"
)

type BoltUserStorage struct {
	db *bolt.DB
}

func NewBoltUserStorage(path string) (*BoltUserStorage, error) {
	var err error
	var db *bolt.DB
	db, err = bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	// init buckets
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte(ADDRESS_CATEGORY))
		tx.CreateBucket([]byte(ADDRESS_ID))
		tx.CreateBucket([]byte(ID_ADDRESSES))
		tx.CreateBucket([]byte(ADDRESS_TIME))
		tx.CreateBucket([]byte(PENDING_ADDRESSES))
		tx.CreateBucket([]byte(CATLOG_PROCESSOR_STATE))
		return nil
	})
	storage := &BoltUserStorage{db}
	return storage, nil
}

func (self *BoltUserStorage) SetLastProcessedCatLogTimepoint(timepoint uint64) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CATLOG_PROCESSOR_STATE))
		err = b.Put([]byte("last_timepoint"), uint64ToBytes(timepoint))
		return err
	})
	return err
}

func (self *BoltUserStorage) GetLastProcessedCatLogTimepoint() (uint64, error) {
	var result uint64
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CATLOG_PROCESSOR_STATE))
		result = bytesToUint64(b.Get([]byte("last_timepoint")))
		return nil
	})
	return result, nil
}

func (self *BoltUserStorage) UpdateAddressCategory(address string, cat string) error {
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		// map address to category
		b := tx.Bucket([]byte(ADDRESS_CATEGORY))
		addrBytes := []byte(strings.ToLower(address))
		b.Put(addrBytes, []byte(strings.ToLower(cat)))
		// get the user of it
		b = tx.Bucket([]byte(ADDRESS_ID))
		user := b.Get(addrBytes)
		if len(user) == 0 {
			// if the user doesn't exist, we set the user to its address
			user = addrBytes
		}
		// add address to its user addresses
		b = tx.Bucket([]byte(ID_ADDRESSES))
		b, err = b.CreateBucketIfNotExists(user)
		if err != nil {
			return err
		}
		b.Put(addrBytes, []byte{1})
		// add user to map
		b = tx.Bucket([]byte(ADDRESS_ID))
		b.Put(addrBytes, user)
		// remove address from pending list
		b = tx.Bucket([]byte(PENDING_ADDRESSES))
		b.Delete(addrBytes)
		return nil
	})
	return err
}

func (self *BoltUserStorage) UpdateUserAddresses(user string, addrs []string, timestamps []uint64) error {
	user = strings.ToLower(user)
	addresses := []string{}
	for _, addr := range addrs {
		addresses = append(addresses, strings.ToLower(addr))
	}
	var err error
	self.db.Update(func(tx *bolt.Tx) error {
		timeBucket := tx.Bucket([]byte(ADDRESS_TIME))
		for _, address := range addresses {
			// get temp user identity
			b := tx.Bucket([]byte(ADDRESS_ID))
			oldID := b.Get([]byte(address))
			// remove the addresses bucket assocciated to this temp user
			b = tx.Bucket([]byte(ID_ADDRESSES))
			b.DeleteBucket(oldID)
			timeBucket.Delete([]byte(address))
			// update user to each address => user
			b = tx.Bucket([]byte(ADDRESS_ID))
			b.Put([]byte(address), []byte(user))
		}
		// remove old addresses from pending bucket
		pendingBk := tx.Bucket([]byte(PENDING_ADDRESSES))
		oldAddrs, _, err := self.GetAddressesOfUser(user)
		if err != nil {
			return err
		}
		for _, oldAddr := range oldAddrs {
			pendingBk.Delete([]byte(oldAddr))
		}
		// update addresses bucket for real user
		// add new addresses to pending bucket
		b := tx.Bucket([]byte(ID_ADDRESSES))
		b, err = b.CreateBucketIfNotExists([]byte(user))
		if err != nil {
			return err
		}
		catBk := tx.Bucket([]byte(ADDRESS_CATEGORY))
		for i, address := range addresses {
			b.Put([]byte(address), []byte{1})
			cat := catBk.Get([]byte(address))
			if string(cat) != KYC_CATEGORY {
				pendingBk.Put([]byte(address), []byte{1})
			}
			log.Printf("storing timestamp for %s - %d", address, timestamps[i])
			timeBucket.Put([]byte(address), uint64ToBytes(timestamps[i]))
		}
		return nil
	})
	return err
}

// returns lowercased category of an address
func (self *BoltUserStorage) GetCategory(addr string) (string, error) {
	addr = strings.ToLower(addr)
	var err error
	var result string
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ADDRESS_CATEGORY))
		cat := b.Get([]byte(addr))
		result = string(cat)
		return nil
	})
	return result, err
}

func (self *BoltUserStorage) GetAddressesOfUser(user string) ([]string, []uint64, error) {
	var err error
	result := []string{}
	timestamps := []uint64{}
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ID_ADDRESSES))
		timeBucket := tx.Bucket([]byte(ADDRESS_TIME))
		userBucket := b.Bucket([]byte(user))
		if userBucket != nil {
			userBucket.ForEach(func(k, v []byte) error {
				addr := string(k)
				result = append(result, addr)
				timestamps = append(timestamps, bytesToUint64(timeBucket.Get(k)))
				return nil
			})
		}
		return nil
	})
	return result, timestamps, err
}

// returns lowercased user identity of the address
func (self *BoltUserStorage) GetUserOfAddress(addr string) (string, uint64, error) {
	addr = strings.ToLower(addr)
	var err error
	var result string
	var timestamp uint64
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ADDRESS_ID))
		timeBucket := tx.Bucket([]byte(ADDRESS_TIME))
		id := b.Get([]byte(addr))
		result = string(id)
		timestamp = bytesToUint64(timeBucket.Get([]byte(addr)))
		return nil
	})
	return result, timestamp, err
}

// returns all of addresses that's not pushed to the chain
// for kyced category
func (self *BoltUserStorage) GetPendingAddresses() ([]string, error) {
	var err error
	result := []string{}
	self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PENDING_ADDRESSES))
		b.ForEach(func(k, v []byte) error {
			result = append(result, string(k))
			return nil
		})
		return nil
	})
	return result, err
}
