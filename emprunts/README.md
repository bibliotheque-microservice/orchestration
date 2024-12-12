# Service_des_emprunts
L'application de gestion d'une bibliothèque


## API Reference

#### PUT emprunts : rendre un livre

```http
  PUT /v1

```

```json
  {
    "empruntId":5, 
    "returned" : true
}
```


#### POST emprunts : emprunter un livre

```http
  POST /v1/emprunt

```

```json
{
    "bookId" : 1,
    "userId": 2
}
```


## MQTT Reference
#### Etre notifié quand un livre est rendu

- queue : "emprunts_finished_queue"
 
- routing key : "emprunts.v1.finished"

example de rendu 

```json
  {
    "CreatedAt": "2024-12-02T12:10:44.741317Z",
    "DateEmprunt": "2024-11-10T10:00:00Z",
    "DateRetourEffectif": "2024-12-02T21:19:02.087307Z",
    "DateRetourPrevu": "2024-11-25T10:00:00Z",
    "DeletedAt": null,
    "IDEmprunt": 5,
    "LivreID": 105,
    "UpdatedAt": "2024-12-02T21:19:02.099121Z",
    "UtilisateurID": 1
}
```


#### Etre notifié quand un emprunt est pour la première fois en retard et donc a une sanction financiète

- queue : "user_penalties_queue"
 
- routing key : "user.v1.penalities.new"

example de rendu 

```json
{
  "penalityId": 3,
  "empruntId": 6,
  "amount": 60.2,
  "userId": 8,
  "created_at": "2024-12-04T10:23:00.016773Z",
  "updated_at": "2024-12-04T22:42:00.011706844Z"
}
```


#### Etre notifié quand un emprunt est à nouveau en retard (un jour de plus), la santion financière est donc mis à jour

- queue : "user_penalties_queue"
 
- routing key : "user.v1.penalities.updated"

example de rendu 

```json
{
  "penalityId": 3,
  "empruntId": 6,
  "amount": 60.2,
  "userId": 8,
  "created_at": "2024-12-04T10:23:00.016773Z",
  "updated_at": "2024-12-04T22:42:00.011706844Z"
}
```



#### Etre notifié quand un nouvel emprunt est crée, le livre n'est donc plus disponible

- queue : "emprunts_created_queue"
 
- routing key : "emprunts.v1.created"

example de rendu 

```json
  {
    "disponible":false,
    "idUtilisateur":2,
    "livreId":1
  }
```