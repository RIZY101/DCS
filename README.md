# DCS
Distributed Cloud Storage
The original inspiration for this project came from an idea for my (Richard Zins) undergraduate capstone project. You may read the document CapstoneProjectProposalRichardZins.pdf in order to gain a better understanding of the purpose, vision, and technical architecture of this project.
## Installation
Plese use the commands bellow to create your TLS keys (run in src directory). Please note if you are connecting to my DCS network you only need to create the key.
```
 >openssl req -new -nodes -x509 -out server.pem -keyout server.key -days 365

 >openssl req -new -nodes -x509 -out client.pem -keyout client.key -days 365
 ```
