# goTest

 Our application is to create two Microservices in backend to securely store data in a file format and allow the user to read and update when required.
 Application supports storing the data in both CSV and XML file format.
 
 Application consists of two services serviceOne and serviceTwo

 ServiceOne is consumer-facing service which will accept request from the user.It exposes a REST-API for achieving this.(when deployed it will use localhost:8000)
 
 Service one and Servicetwo will communicate via gRPC on port 9000.They exchange data using protobubff.
 Based on file type recieved service two will either send/recieve/update data from csv or xml file.
 
 Ensure two install go and protobuff in the host before running the application
