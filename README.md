
# Usefull files : 
server.go    
client.go  
Database.txt

  
# Where are we ? 
Our client is made of a loop that just listens to the server, and answers if it's a question  
Our server accepts the client connection  

First it finds the connected user the best match within all the other users by comparing their music list in the file "Database.txt"  
Then it checks if the connected user has any new subscribers  
If so, it asks if the user is willing the share his music  
In that case, in the file 'utilisateurs.txt', in witch one line equals one client equals one liste of musics, we take the musics from one user, to put it in the list of music of another one   

Then, any client cans start folowing any other one.  
  
# What to do next :
- in the sendData function, we need to find the reason why there still is a line break when writing in the text file  
