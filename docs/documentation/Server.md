# Danger Dodgers Server

This repository contains the server sided code for the danger dodgers project.

## Project Structure

The project is divided according to the following essential structure:

```
DangerDodgersServer
│
└───scripts
└───app
│   └───algo
│   └───app
│   │   └───database
│   │   └───jobs
│   │   └───routers
│   │   └───utils
│   │   │
│   │   │ main.py
│   │
│   │ Dockerfile
│
│   docker-compose
│   requirements.txt
│   .env
```

- Scripts is a directory containing all associated scripts to create any resources for the service. This may include sql scripts to set up the database, or other bash scripts to set up the environment.
- App holds the source code for the application, as well as the scripts necessary to build and run the service. 
- The app child directory within app contains the source code for the back-end application. All source code should remain inside this directory.
- The database directory holds helper functions for connecting and disconnecting with the database.
- The jobs directory should hold any long running jobs which must be processed asynchronously. Examples of long running jobs would be processing a compute intensive algorithm, or saving a large amount of data into the database.
- Routers holds the main CRUD implementations for the service, with each CRUD object in its own python file. If a new CRUD object is introduced to the service, then a new file should be created to house these new routes. Anytime a new set of routes is created, such routes must be registered in the entrypoint file.
- Utils holds any other helper functions which can be abstracted out of the router implementations. This would include rate limiting, user auth, etc. which is used throughout the application.
- Main.py is the entrypoint for the application. It is best to keep logic out of this file and instead implement all route handling in the routers.
- Dockerfile describes the build process for the service, and the dockerfile will be invoked as a part of the docker-compose build. 
- Docker-compose describes how to bring all services up. In the future, as more services are being introduced to the application, this file should contain to bring up the entire back end with a single command.
- Requirements.txt describes the dependencies of the service, and such requirements will be installed as a part of the build process.
- The .env specifies any environmental variables used by the service. 

There are other files contained within the repo, which are for the large part auto generated files created as a part of the build process or docker. 