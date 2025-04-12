# CookieFarm — Stack Tecnologico Definitivo

---
## Server
### 🛠 Backend

### ✨ Framework: [Fiber](https://gofiber.io/)

- Basato su `fasthttp`, garantisce performance eccellenti e bassa latenza.
- Routing veloce e semplice (stile Express.js).    
- Ottimo supporto per middleware, gestione JSON e file statici
\
#### 💾 Database: SQLite

#### 📲 API REST

- Endpoints per ricevere flag (`POST /api/submit`), ottenere statistiche, log, e flags.
- Possibilità di aggiungere WebSocket/SSE per aggiornamenti realtime.

---

### 💻 Frontend

#### 📈 Framework: Nuxt

- Veloce, moderno e leggero grazie a Vite.
- UI basata su componenti, facilmente estendibile.

#### ✨ Styling: TailwindCSS

- Design minimale e responsivo.
- Componenti riutilizzabili, adattabili per dashboard, tabelle, grafici.
#### 🌟 Componenti UI: shadcn/ui

- Collezione di componenti moderni per React con design elegante.
- Perfetti per modali, tab, card, input, selettori, notifiche.

---

## 🛡️ Client (Attacker Bot)

### ✨ Linguaggio: Golang

- Ogni client esegue gli exploit localmente e invia flag al server.
- Concorrenza gestita con goroutine.

---

### 🤝 Dev Tools

- Live reload (in sviluppo): Vite + esecuzione Go hot-reload con `air`.
- Dockerfile per buildare tutto in immagine self-contained.

---

## 🔧 Architettura finale

```
[React (Vite)] → Build → Embedded in Go
     ↳ Servito da Fiber come static asset
     ↳ Comunica con API REST (Fiber)

[Client Bot (Go)] → Esegue exploit → Invia flag → API Server

[Server Fiber]
  ↳ Riceve flag
  ↳ Scrive su SQLite
  ↳ Serve dashboard + dati

[Optional: Redis] → Cache deduplicazione + rate limit
```

---

## 🔄 Obiettivo

Un sistema leggero, performante, modulare e facile da distribuire, che consenta l’automazione delle sottomissioni e il monitoraggio avanzato di exploit e flag in tempo reale.