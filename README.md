
# Usefull files : 
server.go  
user1.go  
user2.go  
music1.go  
music2.go  

# Where are we ? 
user1 and user2 can both reach the server simultaneously  
user1 reads music1 and sends it to user2  
user2 reads music2 and sends it to user1  
when the exchange is done, they are both disconnected  

# What to do next :
start to read our real music files  
create our realtion tree (with a truce matrice)  
make our server run the "exchange" function as soon as a new relation is created between two users  
