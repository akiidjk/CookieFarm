# CookieFarm

## Roadmap

![Roadmap](images/roadmap.png)

## Schema

![Schema](images/schema.png)

# Versions

- Bun: 1.2.9
- Go: 1.24.2
- Docker: 28.0.4


# Go code quality

- Grade .......... A+ 100.0%
- Files ................. 28
- Issues ................. 0
- gofmt ............... 100%
- go_vet .............. 100%
- gocyclo ............. 100%
- ineffassign ......... 100%
- license ............. 100%
- misspell ............ 100%


# NOTE

- release 1.1 è fixing
- release 1.2 è feature delle stat e view
- release 1.3 è feature della cli



# IDEE

## RELEASE 1.0
- [x] Setup project
  - [x] Setup client go project (makefile,logging) - @akiidjk
  - [x] Setup server go project (fiber,air,makefile,logging) - @akiidjk
  - [x] Setup frontend nuxt project (with shadcnui, tailwindcss, typescript, eslint, alias path) - @suga
  - [x] Docker base config - @akiidjk
- [x] Setup repository
  - [x] Setup .gitignore - @akiidjk @suga
  - [x] Setup security settings of repository (branch protection, code scanning, pull request review, code review, pull request approval) - @suga
- [x] Configurazione da file/shitcurl (json parsato)
- [x] Numero configurabile di thread nelle coroutine
- [x] Upgrade FE
- [x] Aggiungere Docs (Codice e user)
- [x] Rifare log e cli di server e client
- [x] Ottimizzazione exploiter per gestire servizi di macchine down etc...
- [x] Fare il README.md
- [x] Aggiungere un bottone sulla table per aggiornarla senza ricaricare la pagina

## RELEASE 1.1
- [ ] Permettere all'utente di stampare la flag (e la merda che vuole) in stdout senza dover restituire obbligatoriamente la flag (by Matte) @Matte
- [x] api/button to send flags, instead of waiting the timer @akiidjk
- [x] api/button to delete/remove the flag queue @akiidjk
- [ ] Filtri, Search flag, Sort @suga
- [x] Auto Reload @akiidjk
- [ ] Compatibilità windows/macos @Dabi
- [ ] Config Reloader (hot realod of config file || button to reload) @akiidjk
- [ ] Tutorial nella dashboard @suga
- [ ] Possibilità di scaricare dal server il client @suga -> dependes on (Tutorial nella dashboard @suga)
- [ ] Non mettere 5 threads di default ma mettere un numero adeguato in base al numero di team presi della config e mettere un upper-limit @suga
- [x] Non fare la batch per la stampa delle flag ma alla conclusione di ogni exploit stampare le flag prese direttamente senza aspettare l'attesa degi altri exploit !!!! PRIORITY
- [X] Implementazione websocket da client a server !!!! PRIORITY
- [x] Inserire la porta non nell'exploit ma come argomento dell'exploiter
- [x] Fornire il numero del team e il nome del servizio nella funzione

## RELEASE 1.2
- [ ] Aggiungere un display che misura ram e cpu del processo (cli - client, web - server)
- [ ] Aggiungere numero client attacanti realtime displayato sulla dashboard,
- [ ] TTL per le flag senza condizione statistiche
- [ ] Completed cli (create template,RealTime consumo di risorse di tutti i client e boh altre info,flag che sono state inviate al server)
- [ ] Exploit Manager che runna più di un exploit (by Matte)

## RELEASE 2.0 (Cyberchallenge update)
- [ ] Calcolatore della SLA
- [ ] Simulatore dell'andamento della gara (active learning)
- [ ] Auto Flag

## RELEASE BOH SI QUANDO ABBIAMO TEMPO
- [ ] Sostituire le richieste in GO con `requests.h`


## Final test

- Test con infra reale
- Team separato
- Exploit dai writeup
