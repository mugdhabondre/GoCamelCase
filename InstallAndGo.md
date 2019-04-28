# installAndGo

Let's see how to get started with goCamelCase!

If you do not have Go installed, [this](https://golang.org/doc/install) is a good place to start.<br />
You can set $GOROOT as an environment variable to your $HOME/go path, as:

```bash
export GOROOT=$HOME/go
echo GOROOT
```  

When you are all set to _Go_, lets download the github repo. 
Type the following commands in bash.

```bash
$ cd $GOROOT/src/
$ go get github.com/mugdhabondre/goCamelCase
$ cd github.com/mugdhabondre/goCamelCase
```

This will install the repo at $GOROOT/src/github.com/mugdhabondre/goCamelCase on your computer.<br />

## Setting up Oxford Dictionary API credentials

The application requires Oxford Dictionaries API credentials. You can create an account on [Oxford Dictionaries website](https://developer.oxforddictionaries.com/).<br />
Create a new file ./credentials.json in the current directory and paste these credentials - app_id and app_key in json format. A sample can be found [here](https://github.com/mugdhabondre/goCamelCase/example.json)

## Lets build it

The repository contains a [build script](https://github.com/mugdhabondre/goCamelCase/build.sh). Run it in bash as follows to build all the src files:

```bash
$ ./build.sh
```

This will create a new _build_ folder and executable ./build/camelcaseapp.

## Deploy
Lets deploy the service on localhost. 
Run the following command in. bash:

```bash
$ build/camelcaseapp
Server started .....
Listen serve for port...  3006
```

By default the http server listens at port number 3006. <br/>
You can go to http://localhost:3006 to view the homepage of the application. 

Lets try requesting for camelCase on a string. Enter the following url:<br />
http://localhost:3006/camelcase/gameofthrones

You will see an output like:
```bash
Input:gameofthrones
Result:gameOfThrones
```

## Deploy in Docker

Now lets try to deploy our solution in a container, so that we can deploy the container on Azure.

[Create](https://hub.docker.com/) a docker hub account if you dont have one. My username on docker hub is _mugdhab_ and gocalecase is the image that I want to build.<br/>
A Dockerfile is included in this repo which creates an image from Alpine Linux.<br/>
Build a container image as follows and try to run it locally:

```bash
$ docker build -t mugdhab/gocamelcase .
$ docker run -p 5005:80 --env PORT=80 mugdhab/gocamelcase:latest 
```

You can access the application at http://localhost:5005/ .

Lets publish the image. on docker hub:

```bash
$ docker login
$ docker push mugdhab/gocamelcase:latest
``` 

Now we can access this image in the docker hub repo and use it for our Azure app!

## Deploy in Azure

I followed [this](https://medium.com/@durgaprasadbudhwani/step-by-step-guide-to-deploy-golang-application-on-azure-web-app-46ba3befb4c0) article to deploy this application in Azure.

Sign in to Azure Portal and open cloud shell button to open a cloud PowerShell. 

Type the following commands to create the application:
```bash
$ az group create --name AzureGoRG --location "South Central US"
$ az appservice plan create --name AzureGoSP --resource-group AzureGoRG --sku S1 --is-linux
$ az webapp create --resource-group AzureGoRG --plan AzureGoSP  --name camelcasegenerator --deployment-container-image-name mugdhab/gocamelcase:latest
``` 

Now the application will. be available @ https://camelcasegenerator.azurewebsites.net/. Yay!

