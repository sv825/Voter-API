package main

import (
	"testing"
	"time"
)

func TestVoterList(t *testing.T) {
	vl := NewVoterList()
	vl.AddVoter(Voter{
		VoterID:   1,
		FirstName: "John",
		LastName:  "Doe",
	})
	voters := vl.GetVoters()
	if len(voters) != 1 {
		t.Errorf("Expected 1 voter, got %d", len(voters))
	}
	voter, ok := vl.GetVoter(1)
	if !ok {
		t.Error("Expected to get voter with ID 1")
	}
	if voter.FirstName != "John" {
		t.Errorf("Expected first name to be John, got %s", voter.FirstName)
	}
	vp := voterPoll{
		PollID:   1,
		VoteDate: time.Now(),
	}
	vl.AddVoterPoll(1, vp)
	voter, ok = vl.GetVoter(1)
	if !ok {
		t.Error("Expected to get voter with ID 1")
	}
	if len(voter.VoteHistory) != 1 {
		t.Errorf("Expected 1 vote history, got %d", len(voter.VoteHistory))
	}
}

package main

import (
	"testing"
	"time"
)

func TestVoterPoll(t *testing.T) {
	vp := voterPoll{
		PollID:   1,
		VoteDate: time.Now(),
	}
	if vp.PollID != 1 {
		t.Errorf("Expected PollID to be 1, got %d", vp.PollID)
	}
}

func TestVoter(t *testing.T) {
	v := Voter{
		VoterID:   1,
		FirstName: "John",
		LastName:  "Doe",
	}
	if v.VoterID != 1 {
		t.Errorf("Expected VoterID to be 1, got %d", v.VoterID)
	}
	if v.FirstName != "John" {
		t.Errorf("Expected FirstName to be John, got %s", v.FirstName)
	}
	if v.LastName != "Doe" {
		t.Errorf("Expected LastName to be Doe, got %s", v.LastName)
	}
	vp := voterPoll{
		PollID:   1,
		VoteDate: time.Now(),
	}
	v.VoteHistory = append(v.VoteHistory, vp)
	if len(v.VoteHistory) != 1 {
		t.Errorf("Expected VoteHistory to have length 1, got %d", len(v.VoteHistory))
	}
}


