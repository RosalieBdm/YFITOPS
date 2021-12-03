
# Usefull files : 
server.go    
client.go  
utilisateurs.txt

  
# Where are we ? 
Our client is made of a loop that just listens to the server, and answers if it's a question  
Our server accepts the client connection  

First it checks if the connected user has any new subscribers  
If so, it asks if the user is willing the share his music  
In that case, in the file 'utilisateurs.txt', in witch one line equals one client equals one liste of musics, we take the musics from one user, to put it in the list of music of another one   

Then, any client cans start folowing any other one.  
  
# What to do next :
create our relation tree (with a truce matrice)  
--> if a client subscribes to another one, he doesn't only get his musics list once, but it actually evolves to stay the same even he the other changes his list  
  
