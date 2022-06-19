# Title 
Remove learner role from paxos
# Date 
12/06/2022 
# Status 
Accepted 
# Context
We have decided to remove learner role from our paxos implemention.

In paxos, learner is role that don't have ability to propose but only accept proposal from others.
We designed our system to be able to regist and execute task from any node, so learner role is not needed.

# Decision 
Since there are only two of us developing, we'll make it easy by cut off unnecessary entities. 

# Consequences
Less realistic to the paxos problem