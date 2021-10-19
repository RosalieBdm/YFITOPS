# Suivi du projet
*Le but de ce fichier est de suivre l'évolution de notre projet est d'en garder une trace écrite notamment pour les bilans*  

## Séance 1 
mail prof : pierre.francois@insa-lyon.fr
Projet :
 - groupe de 3 de preference mixtes  
 - A la fin de chaque cours bilan, même minime, au prof par mail avec objet specifique : [GO] GR2-<numProjet> etc  

Objectif de cette séance :   
    1. Installation GO  
    2. Formation du groupe  
    3. Programme qui ouvre un fichier, lire le fichier et afficher à l'écran les derniers mots de chaque ligne   
    4. Créer un repo git   
    5. Préparer le mail    

## Notre groupe : GR2-8 
Piste de recherche :   
 - Traitement d'image   
 - Compression/amélioration qualité d'image    

Premier travail : tranformer une image en noir et blanc    
Réussite ! &rarr; code dispo dans *image_bw/main.go*    

A voir pour go routine &rarr; systeme de canal    
Sinon plusieurs go routine qui ecrive dans le meme fichier mais il faut verifier si il est thread safe -> mutex     

## Séance 2 (autonomie)
On test notre programe avec/sans go routines pour une image en 6000x4000 donc plutot lourde.   
Sans go routines : 2,75 secondes   
Avec go routine : envrion 10 secondes    
Conclusion : il y a un problème dans le code avec des go routines il va falloir trouver un moyen d'optimiser la chose   

## Séance 3 du 08/10/2021
Début de reflexion :  
On est passé dans un premier temps sur un buffered channel -> gain de 2 secondes  
On a ensuite essayé de créer un autre programme fonctionnant avec un système de mutex qui va directement écrire chaque pixel dans le fichier final on se retrouve avec un temps d'execution avoisinant les 4 secondes pour 4 subdivisions. Ce n'est toujours pas mieux que le programe sans go routines mais on s'améliore.  
On a enfin réussi à rendre les go routines plus efficace !  
Solutions ? Dans la go routine on utilisait un mutex qui bloquait à chaque écriture de pixel or cela n'est pas efficaace car il faudrait faire un blocage si jamais deux go routines tentent d'écrire au meme endroit ce qui n'est pas le cas pour nous : chaque go routines écrit dans des pixels différents donc pas de risque de collision.  
Avec ces nouvelles modifications on obtient les résultats suivants (on ne compte que le temps de traitements de l'image on exclue le temps de lecture de l'input et de l'ecriture de l'output)  
 - Sans go routine envrion 2 secondes   
 - Avec go routine 400-410 ms avec 8 go routines  
Donc gros gain apporté par les go routines.

## Séance 4 du 14/10/2021
Depuis la dernière fois : 
 - Ajout de commentaire sur le code 
 - Ajout de PFR en tant que collaborateur du git
 - Décision d'avancer à la partie serveur car go routine et traitement d'image acquis

Début de séance : introduction à TCP      
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
Mais attention cette méthode (break au moment de l'erreur signifiant la fin du fichier) ne correspond pas à tout les types de problèmes. En fonction de ce qu'on envoie il fat faire comprendre au serveur qu'on a finit l'envoie et qu'on veut lancer le traitement. Par exemple pour nos images.   
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
On va alors adapter cela à notre code.      
Pour tranférer image &rarr; envoyer un go object pour se simplifier la vie. Le seul "inconvénient" est qu'il faut que le serveur et le client soient en golang.   
Pour se familiariser avec le transfer de structure on s'est alors basé sur le code disponible [ici](https://gist.github.com/MilosSimic/ae7fe8d70866e89dbd6e84d86dc8d8d5) qui nous a permis de comprendre comment envoyer une structure assez simple. Notre but va alors de bien comprendre et tranformer ce code afin de pouvoir tranférer des images sur le server TCP.   
