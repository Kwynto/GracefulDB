# GracefulDB
Fast, Simple and Secure. 
This is a DBMS for professionals and extreme loads. 

**_This repository is under development._**

**Usage**

> go build -o gdb
> mkdir ./data
> ./gdb

**Testing**

Run tests:
> go test ./... -v

Run tests showing code coverage:
> go test ./... -cover -v

You can view code coverage in detail in your web browser.  
To do this, you need to sequentially execute two commands in the console:
> go test ./... -coverprofile="coverage.out" -v  
> go tool cover -html="coverage.out"

## About the author

The author of the project is Constantine Zavezeon (Kwynto).  
You can contact the author by e-mail: kwynto@mail.ru  
The author accepts proposals for participation in open source projects,  
as well as willing to accept job offers.
If you want to offer me a job, then first I ask you to read [this](https://github.com/Kwynto/Kwynto/blob/main/offer.md).
