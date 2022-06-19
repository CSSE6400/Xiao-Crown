# Title 
Limit max number of nodes in cluster to 5
# Date 
15/06/2022 
# Status 
Accepted 
# Context
Ideally, a distributed cluster is able to change it's configuration dynamically (on-the-fly) without temporarily shutdown the service. This feature would bring more reliability and scalability to the System.

However, we underestimate the complexity of this feature. It is unrealistic to assume all the nodes would switch to the new configuration at same time. That is during the configuration change, it is inevitably some nodes are running with new configuration and others are still in the old mode. This will lead to split-brain problem, made our system unusable.

# Decision 
In order to keep the availability of our system, we set a max nodes limitation in our system.

# Consequences
Less flexibility and scalability.