package repository

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"

	proto "github.com/jakewright/home-automation/service.device-registry/proto"
	"github.com/jinzhu/copier"
)

// RoomRepository provides access to the underlying storage layer
type RoomRepository struct {
	// ConfigFilename is the path to the room config file
	ConfigFilename string

	// ReloadInterval is the amount of time to wait before reading from disk again
	ReloadInterval time.Duration

	rooms    []*proto.Room
	reloaded time.Time
	lock     sync.RWMutex
}

// FindAll returns all rooms
func (r *RoomRepository) FindAll() ([]*proto.Room, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	var rooms []*proto.Room
	for _, room := range r.rooms {
		out := &proto.Room{}
		if err := copier.Copy(out, room); err != nil {
			return nil, err
		}
		rooms = append(rooms, out)
	}

	return rooms, nil
}

// Find returns a room by ID
func (r *RoomRepository) Find(id string) (*proto.Room, error) {
	if err := r.reload(); err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, room := range r.rooms {
		if room.ID == id {
			out := &proto.Room{}
			if err := copier.Copy(out, room); err != nil {
				return nil, err
			}

			return out, nil
		}
	}

	return nil, nil
}

func (r *RoomRepository) reload() error {
	// Skip if we've recently reloaded
	if r.reloaded.Add(r.ReloadInterval).After(time.Now()) {
		return nil
	}

	data, err := ioutil.ReadFile(r.ConfigFilename)
	if err != nil {
		return err
	}

	var cfg struct {
		Rooms []*proto.Room `json:"rooms"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	r.rooms = cfg.Rooms

	r.reloaded = time.Now()
	return nil
}
