Guide

- First run an IPFS daemon and the tedermint node 

- Create the private key for the gonverment
$ ./client g --filename=gon.json
The generate was successful

$ cat gon.json 
{"PublicKey":"08011220d6f9ba28873e213cc8715c7d6cdac7898059b264e03ef127318d97c4e8d88d7b","PrivateKey":"CAESYJ2Admpxx4k57L+k7ef8k5YkSrtDVsOsMuw3nu+RKIDB1vm6KIc+ITzIcVx9bNrHiYBZsmTgPvEnMY2XxOjYjXvW+boohz4hPMhxXH1s2seJgFmyZOA+8ScxjZfE6NiNew=="}

- Start the tendervoting's server by adding public key for the gonverment 
$ ./server -gonverment=08011220d6f9ba28873e213cc8715c7d6cdac7898059b264e03ef127318d97c4e8d88d7b
I[06-17|16:42:03.481] Starting ABCIServer                          module=abci-server impl=ABCIServer
I[06-17|16:42:03.504] Waiting for new connection...                module=abci-server 

- For simplicity will add the gonverment as the voter for this election, like a good monarch. 
  Only the public keys in the list can vote.
$ ./client ce --key=gon.json --voters=08011220d6f9ba28873e213cc8715c7d6cdac7898059b264e03ef127318d97c4e8d88d7b
The election submitted with ID 09202777-6d10-49e1-b310-1843a2731af1

- We will query the current elections
$ ./client e
Election ID: 09202777-6d10-49e1-b310-1843a2731af1
Latest: true
Number of voters: 1

- If we add more than one election, then we will need the latest election's ID to create a poll.
  Only with the latest election we can create a poll.
$ ./client le
Election ID: 09202777-6d10-49e1-b310-1843a2731af1
Latest: true
Number of voters: 1

- Now we will create a directory for the poll and the poll.json
mkdir /tmp/examplePoll
echo '{"Description":"Is tendervoting cool ?","Choices":{"y":"Yes","n":"No"}}' > /tmp/examplePoll/poll.json

- We will submit the directory to the IPFS
$ ipfs add -r /tmp/examplePoll/
added QmatKDJ8Nf1dQaWzZi24GsBcTLsoy2X6hxKEAi8ryALgfD examplePoll/poll.json
added QmPdy89ZQt4c6EWECMPibPfjFHhe235XKHZNDiAZD5x5tH examplePoll

- We will add the poll to the blockchain, by using the last IPFS hash of the directory
$ ./client cp --key=gon.json --hash=QmPdy89ZQt4c6EWECMPibPfjFHhe235XKHZNDiAZD5x5tH --election=09202777-6d10-49e1-b310-1843a2731af1
The poll submitted

- Now we will query the list of polls
$ ./client p
Poll's Hash: QmPdy89ZQt4c6EWECMPibPfjFHhe235XKHZNDiAZD5x5tH
Latest: true

- If we add more than one poll, then we will need the latest poll's hash to vote.
  Only with on latest poll we can vote.
$ ./client lp
Poll's Hash: QmPdy89ZQt4c6EWECMPibPfjFHhe235XKHZNDiAZD5x5tH
Latest: true

- Now we are going to vote yes, because we do like tendervoting !
$ ./client v --hash=QmPdy89ZQt4c6EWECMPibPfjFHhe235XKHZNDiAZD5x5tH --choice=y --key=gon.json
The vote submitted

- Now we can query the results from the poll.
$ ./client r --hash=QmPdy89ZQt4c6EWECMPibPfjFHhe235XKHZNDiAZD5x5tH
Votes for choice 'n': 0
Votes for choice 'y': 1
Number of voters: 1