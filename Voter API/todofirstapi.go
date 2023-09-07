package main

import (
	"fmt"
	"time"
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
	Voters map[uint]Voter //A map of VoterIDs as keys and Voter structs as values
}

// constructor for VoterList struct
func NewVoterList() *VoterList {
	return &VoterList{
		Voters: make(map[uint]Voter),
	}
}

//Add receivers to any structs you want, but at the minimum you should add the API behavior to the
//VoterList struct as its managing the collection of voters.  Also dont forget in the constructor
//that you need to make the map before you can use it - make map[uint]Voter

// Get all voter resources including all voter history for each voter (note we will discuss the concept of "paging" later, for now you can ignore)
func (vl *VoterList) GetVoters() []Voter {
	voters := make([]Voter, 0, len(vl.Voters))
	for _, v := range vl.Voters {
		voters = append(voters, v)
	}
	return voters
}

// Get a single voter resource with voterID=:id including their entire voting history.
func (vl *VoterList) GetVoter(id uint) (Voter, bool) {
	v, ok := vl.Voters[id]
	return v, ok
}

// POST version adds one to the "database"
func (vl *VoterList) AddVoter(v Voter) {
	vl.Voters[v.VoterID] = v
}

// Gets the JUST the voter history for the voter with VoterID = :id
func (vl *VoterList) GetVoterPolls(id uint) ([]voterPoll, bool) {
	v, ok := vl.Voters[id]
	if !ok {
		return nil, false
	}
	return v.VoteHistory, true
}

// Gets JUST the single voter poll data with PollID = :id and VoterID = :id.
func (vl *VoterList) GetVoterPoll(voterID uint, pollID uint) (voterPoll, bool) {
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
	v, ok := vl.Voters[voterID]
	if !ok {
		return false
	}
	v.VoteHistory = append(v.VoteHistory, vp)
	vl.Voters[voterID] = v
	return true
}

// Returns a "health" record indicating that the voter API is functioning properly and some metadata about the API.
func (vl *VoterList) HealthCheck() string {
	return "API is functioning properly"
}

func main() {
	vl := NewVoterList()
	vl.AddVoter(Voter{
		VoterID:   1,
		FirstName: "John",
		LastName:  "Doe",
	})
	fmt.Println(vl.GetVoters())
	fmt.Println(vl.GetVoter(1))
	fmt.Println(vl.GetVoterPolls(1))
	vl.AddVoterPoll(1, voterPoll{
		PollID:   1,
		VoteDate: time.Now(),
	})
	fmt.Println(vl.GetVoterPolls(1))
	fmt.Println(vl.GetVoterPoll(1, 1))
	fmt.Println(vl.HealthCheck())
}
