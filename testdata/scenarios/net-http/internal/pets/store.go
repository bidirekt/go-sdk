package pets

import (
	"sort"
	"sync"
	"time"
)

type Store struct {
	mutex  sync.Mutex
	pets   map[int64]Pet
	nextID int64
}

func NewStore() *Store {
	return &Store{pets: map[int64]Pet{}, nextID: 1}
}

func (this *Store) All() []Pet {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	all := make([]Pet, 0, len(this.pets))
	for _, pet := range this.pets {
		all = append(all, pet)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	return all
}

func (this *Store) Find(id int64) (Pet, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	pet, found := this.pets[id]
	return pet, found
}

func (this *Store) Add(newPet NewPet) Pet {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	pet := Pet{
		ID:        this.nextID,
		Name:      newPet.Name,
		Status:    StatusAvailable,
		Tags:      newPet.Tags,
		Weight:    newPet.Weight,
		BirthDate: newPet.BirthDate,
		CreatedAt: time.Now().UTC(),
	}
	this.pets[pet.ID] = pet
	this.nextID++
	return pet
}

func (this *Store) ChangeStatus(id int64, status Status) (Pet, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	pet, found := this.pets[id]
	if !found {
		return Pet{}, false
	}
	pet.Status = status
	this.pets[id] = pet
	return pet, true
}

func (this *Store) Remove(id int64) bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, found := this.pets[id]; !found {
		return false
	}
	delete(this.pets, id)
	return true
}
