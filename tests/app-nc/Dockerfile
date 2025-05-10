FROM golang:latest

# Installa socat e Go dependencies di base
RUN apt update && apt install -y socat build-essential

# Setta la working directory
WORKDIR /app

# Copia il sorgente Go
COPY . .

# Compila il binario (assumendo che il file si chiami app-nc.go)
RUN go build -o app-nc ./main.go

# Rendi eseguibile e pulisci eventuali file non necessari
RUN chmod +x ./app-nc

# Espone la porta 1337
EXPOSE 1338

# Usa socat per ascoltare su TCP e redirigere all'app
CMD ["socat", "-t", "10", "-T", "10", "TCP-LISTEN:1338,reuseaddr,fork", "EXEC:/app/app-nc"]
