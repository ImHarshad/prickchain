# prickchain
A small POC of simple private Blockchain.

# Setup
After installing and configuring Go from here - https://golang.org/dl/

Grab the following packages:

Spew allows us to view structs and slices cleanly formatted in our console. This is nice to have.

    go get github.com/davecgh/go-spew/spew

Gorilla/mux is a popular package for writing web handlers. We’ll need this.

    go get github.com/gorilla/mux

Godotenv lets us read from a .env file that we keep in the root of our directory so we don’t have to hardcode things like our http ports.

    go get github.com/joho/godotenv

Now in your workspace clone Prickchain.

Fire up your application from terminal using command -

    go run main.go
