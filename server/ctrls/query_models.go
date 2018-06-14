package ctrls

type PollQuery struct {
	PollHash string
}

type ItemPollQuery struct {
	PollQuery
	Latest bool
}

type ListPollQuery []ItemPollQuery

type ElectionQuery struct {
	ID             string
	NumberOfVoters int
}

type ItemElectionQuery struct {
	ElectionQuery
	Latest bool
}

type ListElectionQuery []ItemElectionQuery

type PollVotesQuery struct {
	Choices       map[string]int
	NumberOfVotes int
}
