# TCP Server
Task
Write a TCP Server that loads dataset from a CSV file and provide interface to query the dataset. TCP server will expose port 4040.

Server Communication
User should be able to connect to the server using NetCat nc localhost 4040 command on Linux/Unix based systems.
Once connected to TCP, user should be able communicate with the application by sending queries in JSON format.

User can query data based on two fields: Region and Date.
In response, server will return a list of records that will match query.

