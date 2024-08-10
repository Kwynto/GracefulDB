# GracefulDB
Fast, Simple and Secure. 
This is a DBMS for professionals and extreme loads. 

**_This repository is under development._**

## Using 

The operating instructions are not ready yet.  
You can compile and run a project just like any other GoLang project.  


Download the GracefulDB project to your computer:  
> git clone https://github.com/Kwynto/GracefulDB.git  

or in another way  

Go to the project folder:  
> cd ./GracefulDB  

Compile the project:  
> go build main.go

Start the server  
For Windows:  
> .\main.exe  

For *nix:  
> ./main

**Warning:** Users of Unix systems may need to change the access rights for the executable file.  

The server is running and now you can manage it. To do this, go to the web browser:  
> http://localhost  

or  

> http://localhost:80  

## Testing 

Run tests:  
> go test ./... -v  

Run tests showing code coverage:  
> go test ./... -cover -v  

You can view code coverage in detail in your web browser.  
To do this, you need to sequentially execute two commands in the console:  
> go test ./... -coverprofile="coverage.out" -v  
> go tool cover -html="coverage.out"  


## Thanks for the help

*this section is still empty*  

You can support this project and your name or the name of your company can take its place in our hall of fame. The details are [here](https://github.com/Kwynto/GracefulDB/blob/main/SUPPORT.md).


## About the author 

The author of the project is Constantine Zavezeon (Kwynto).  
You can contact the author by e-mail: kwynto@mail.ru  
The author accepts proposals for participation in open source projects,  
as well as willing to accept job offers.
If you want to offer me a job, then first I ask you to read [this](https://github.com/Kwynto/Kwynto/blob/main/offer.md).
