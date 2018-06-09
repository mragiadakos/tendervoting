Tendervoting

- The voter creates a public key and print it in a QR code. This way the voter can give it to the gonverment.
Library to print QR code  github.com/skip2/go-qrcode/
- The gonverment gets the QR code as png it and uploads it to the system
Library to read the QR code github.com/tuotoo/qrcode
- The gonverment uploads also the poll to the system. 
The poll will be and IPFS directory that will contain all the files about the poll, even HTML/CSS.
However the directory will contain a JSON file, called poll.json.
The poll.json will continue two attributes: Description and Choices.
The gonverment gives the IPFS hash of the directory to the voters.
- The voter gets the IPFS hash to read the description like this:
./client show --poll=< hash >
Description:
bla bla bla
Choices
a) vote for bla
b) vote for bla bla
c) vote for bla bla bla
- Now the voter can vote, however he can vote only once for the specific poll
./client vote --poll=< hash > --choice=< choice id >
The vote was successful 


Delivery
REQUEST for the gonverment
{
    Signature: string
    Data: {
        ID: uuid 
        From: public key as hex
        Voters: array of public keys as hex
        EndTime : time
    }
}
RESPONSE
  Error scenarios:
    - The public key of the gonverment is not in the list

Delivery
REQUEST for the voter
{
    Signature: string
    Data: {
        From: public key as hex 
        PollHash: string 
        Choice: string
        Date: time
    }
}
RESPONSE
  Error scenarios:
    - the voter is not in the system
    - the voter has vote already for the specific PoolHash
    - the voter is not authorized becaused it passed the EndTime


Query
REQUEST
{
    PollHash: string
}

RESPONSE
{
    Votes:{
        Choices: map[string]int // the choices and percentage based on the number of voters
        NumberOfVoters: int
        NumberOfVotes: int
    }
}