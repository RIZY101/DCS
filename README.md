# DCS
Distributed Cloud Storage
The original inspiration for this project came from an idea for my (Richard Zins) undergraduate capstone project. You may read the document CapstoneProjectProposalRichardZins.pdf in order to gain a better understanding of the purpose, vision, and technical architecture of this project. Please note that the goal of this project has moved away from the fully decentralized abilities outlined in the proposal. In place I decided to implement a Cooperative Cloud Storage Network. You can read more about it [here](https://en.wikipedia.org/wiki/Cooperative_storage_cloud)
## Installation
After downloading the repository on multiple machines please make sure you have opened port 6633 on the machine/network your MasterNode is on, and port 6634 for the machine/network your LNode is running on. Once that is done CNode should be able to intiate connections to all parties. Also pleaemake sure you create your own TLS keys instead of useing mine. Mine are only included in this repo for eas of tranfer for myself between multiple machines.
Plese use the commands bellow to create your TLS keys (run in src directory). Please note if you are connecting to my DCS network you only need to create the client key.
```
 >openssl req -new -nodes -x509 -out server.pem -keyout server.key -days 365

 >openssl req -new -nodes -x509 -out client.pem -keyout client.key -days 365
 ```
