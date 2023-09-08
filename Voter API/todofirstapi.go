package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type voterPoll struct {
	PollID   uint
	VoteDate time.Time
}

type Voter struct {
	VoterID     uint
	FirstName   string
	LastName    string
	VoteHistory []voterPoll
}

type VoterList struct {
	sync.RWMutex
	Voters    map[uint]Voter
	bootTime  time.Time
	totalAPI  uint64
	errorAPI  uint64
}

//constructor for VoterList struct
func NewVoterList() *VoterList {
	return &VoterList{
		Voters:   make(map[uint]Voter),
		bootTime: time.Now(),
	}
}

// Get all voter resources including all voter history for each voter (note we will discuss the concept of "paging" later, for now you can ignore)
func (vl *VoterList) GetVoters() []Voter {
	vl.RLock()
	defer vl.RUnlock()
	voters := make([]Voter, 0, len(vl.Voters))
	for _, v := range vl.Voters {
		voters = append(voters, v)
	}
	return voters
}

// Get a single voter resource with voterID=:id including their entire voting history.
func (vl *VoterList) GetVoter(id uint) (Voter, bool) {
	vl.RLock()
	defer vl.RUnlock()
	v, ok := vl.Voters[id]
	return v, ok
}

// POST version adds one to the "database"
func (vl *VoterList) AddVoter(v Voter) {
	vl.Lock()
	defer vl.Unlock()
	vl.Voters[v.VoterID] = v
}

// PUT version updates one in the "database"
func (vl *VoterList) UpdateVoter(id uint, v Voter) bool {
	vl.Lock()
	defer vl.Unlock()
	if _, ok := vl.Voters[id]; !ok {
		return false
	}
	v.VoterID = id
	vl.Voters[id] = v
	return true
}

// DELETE version removes one from the "database"
func (vl *VoterList) DeleteVoter(id uint) bool {
	vl.Lock()
	defer vl.Unlock()
	if _, ok := vl.Voters[id]; !ok {
		return false
	}
	delete(vl.Voters, id)
	return true
}

// Gets the JUST the voter history for the voter with VoterID = :id
func (vl *VoterList) GetVoterPolls(id uint) ([]voterPoll, bool) {
	vl.RLock()
	defer vl.RUnlock()
	v, ok := vl.Voters[id]
	if !ok {
		return nil, false
	}
	return v.VoteHistory, true
}

// Gets JUST the single voter poll data with PollID = :id and VoterID = :id.
func (vl *VoterList) GetVoterPoll(voterID uint, pollID uint) (voterPoll, bool) {
	vl.RLock()
	defer vl.RUnlock()
	v, ok := vl.Voters[voterID]
	if !ok {
		return voterPoll{}, false
	}
	for _, vp := range v.VoteHistory {
		if vp.PollID == pollID {
			return vp, true
		}
	}
	return voterPoll{}, false
}

// POST version adds one to the "database"
func (vl *VoterList) AddVoterPoll(voterID uint, vp voterPoll) bool {
	vl.Lock()
	defer vl.Unlock()
	v, ok := vl.Voters[voterID]
	if !ok {
		return false
	}
	v.VoteHistory = append(v.VoteHistory, vp)
	vl.Voters[voterID] = v
	return true
}

// PUT version updates one in the "database"
func (vl *VoterList) UpdateVoterPoll(voterID uint, pollID uint, vp voterPoll) bool {
	vl.Lock()
	defer vl.Unlock()
	v, ok := vl.Voters[voterID]
	if !ok {
		return false
	}
	found := false
	for i, p := range v.VoteHistory {
		if p.PollID == pollID {
			found = true
			vp.PollID = pollID
			v.VoteHistory[i] = vp
			break
		}
	}
	if !found {
		return false
	}
	vl.Voters[voterID] = v
	return true
}

// DELETE version removes one from the "database"
func (vl *VoterList) DeleteVoterPoll(voterID uint, pollID uint) bool {
	vl.Lock()
	defer vl.Unlock()
	v, ok := vl.Voters[voterID]
	if !ok {
		return false
	}
	found := false
	for i, p := range v.VoteHistory {
		if p.PollID == pollID {
			found = true
			v.VoteHistory = append(v.VoteHistory[:i], v.VoteHistory[i+1:]...)
			break
		}
	}
	if !found {
		return false
	}
	vl.Voters[voterID] = v
	return true
}

// Returns a "health" record indicating that the voter API is functioning properly and some metadata about the API.
func (vl *VoterList) HealthCheck() map[string]interface{} {
	vl.RLock()
	defer vl.RUnlock()
	return map[string]interface{}{
		"uptime":           time.Now().Sub(vl.bootTime).String(),
		"total_api_calls":  vl.totalAPI,
		"total_error_calls": vl.errorAPI,
	}
}

func main() {
	vl := NewVoterList()
	r := mux.NewRouter()
	r.HandleFunc("/voters", func(w http.ResponseWriter, r *http.Request) {
		vl.totalAPI++
		switch r.Method {
		case http.MethodGet:
			voters := vl.GetVoters()
			json.NewEncoder(w).Encode(voters)
		case http.MethodPost:
			var v Voter
			err := json.NewDecoder(r.Body).Decode(&v)
			if err != nil {
				vl.errorAPI++
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			vl.AddVoter(v)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(v)
		default:
			vl.errorAPI++
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodGet, http.MethodPost)

	r.HandleFunc("/voters/{id}", func(w http.ResponseWriter, r *http.Request) {
		vl.totalAPI++
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			vl.errorAPI++
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			voter, ok := vl.GetVoter(uint(id))
			if !ok {
				vl.errorAPI++
				http.Error(w, "voter not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(voter)
		case http.MethodPost:
			var v Voter
			err := json.NewDecoder(r.Body).Decode(&v)
			if err != nil {
				vl.errorAPI++
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			vl.AddVoter(v)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(v)
		case http.MethodPut:
			var v Voter
			err := json.NewDecoder(r.Body).Decode(&v)
			if err != nil {
				vl.errorAPI++
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ok := vl.UpdateVoter(uint(id), v)
			if !ok {
				vl.errorAPI++
				http.Error(w, "voter not found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(v)
		case http.MethodDelete:
			ok := vl.DeleteVoter(uint(id))
			if !ok {
				vl.errorAPI++
				http.Error(w, "voter not found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
