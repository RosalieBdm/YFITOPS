
# Introduction à TCP      
Client &rarr; personne qui décide de se connecter au server       
TCP fiable &rarr; tout paquet envoyé arrive et dans l'ordre    
On va alors discuter de l'implémentation de TCP en go.     
L'application s'enregistre sur l'OS en précisant le port qui la concerne.    
Attention si on est pas root sur sa machine **il faut choisir un port supérieur à 1024** par question de sécurité.    
Donc en go il faut importer le package *net*.
Exemple d'ouverture de serveur TCP :     
```Go
import (
    "net"
    "bufio"
    "fmt"
    "strings"
    "io"
)

func main (){
    ln, err := net.Listen("tcp", ":port") //avec port le numéro du port, ln = listener
    if err != nil { //si jamais on detecte une erreur   
        panic(err)
    }
    connum := 0 //permet de débug en gardant ke nb de connection
    for { //Boucle infinie pour traiter les clients 
        conn, errconn := ln.Accept() //On accepte la connection et on met l'identifiant de la session dans conn
        //Cette ligne bloque le code tant qu'il n'y a pas de connectiom
        if errconn != nil {
            panic(errconn)
        }
        //On prend tout de suite en charge la connection
        go handleConnection(conn, connum)
        connum +=1
    }
}

func handleConnection(connection net.Conn, connum int){
    defer connection.Close() //permet de fermer la connection une fois le code fini !!!! hyper important  
    connReader := bufio.NewReader(connection)
    //Server qui lit des chaînes de caractères et renvoie le dernier mot de chque ligne
    for {
        inputLine, err := connReader.ReadString("\n")
        if err != nil {
            fmt.Printf("problème")
            break //ici on ne panic pas car sinon on tue le serveur alors qu'une erreur va signifier la fin de ligne ou la déconnection d'un client 
        }

        inputLine = strings.TrimSuffix(inputLine, "\n") //TrimSuffix permet de dégager le \n
        splitLine := strings.Split(inputLine, " ") //renvoie un slice en séparant avec le caractère précise, ici l'espace
        returnedString := splitLine[len(splitLine) - 1] //On récupère le dernier mot
        io.WriteString(connection, fmt.Sprintf("%s\n", returnedString))
    }
}
```
Serveur c'est une boucle infinie qui attend la connection d'un client.      
On va paralléliser le traitement des clients avec des go routines.    
Mais attention cette méthode (break au moment de l'erreur signifiant la fin du fichier) ne correspond pas à tout les types de problèmes. En fonction de ce qu'on envoie il faut faire comprendre au serveur qu'on a finit l'envoie et qu'on veut lancer le traitement.   
Regarder : parser , TLV &rarr; Type Length Value    
Regardons maintenant le client :    
```Go
import (
    "fmt"
    "net"
    "os"
)

func main(){
    conn, err := net.Dial("tcp", "127.0.0.1:10000") //Connection sur le port 10000, Rappel : 127.0.0.1 = moi  
    if err != nil {
        os.Exit()
    } else {
        //Traitement 
    }
}
```

# Liens utiles
Go library : https://pkg.go.dev/std

TCP server/client : https://gist.github.com/MilosSimic/ae7fe8d70866e89dbd6e84d86dc8d8d5
POO : https://devopssec.fr/article/programmation-orientee-objet-golang#begin-article-section
Goroutine : https://devopssec.fr/article/goroutines-golang
Gestion des fichiers txt : https://devopssec.fr/article/lire-et-ecrire-dans-un-fichier-golang
Communiquer entre les goroutines : https://devopssec.fr/article/channels-golang
