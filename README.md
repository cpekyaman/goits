# GOITS
**Go Issue Tracking Service** is an issue tracker with backend written in go and front end built on top of Nuxt.js with vuetify as the main component library.

## goits-devtools
Mostly cli utilities for things like code generation.

## goits-server
This is the backend of goits application. 

### How To Run
It uses `make` for building the server, so you need to install make if your system does not have it. Once it is installed, you can run `make build` under `goits-server` to build the executable as `goits/goits.exe`. Then you can start the server by executing `goits server start` from command line. 

### App Structure
- **application** : contains business related packages for the application
- **cli** : code for command line interface of the application
- **config** : code for configuration wrapper / helper
- **framework** : non-business and supporting code such as orm layer
- **main** : just the main entry point of the application
- **server** : contains http server part

## goits-ui
This is the front end of goits application. 
It follows the common directory structure of a nodejs / nuxt webapp.

## bin
Output folder for executable build targets.

## codegen
Root directory for code generation templates and temporary code generation output.

## etc
Configuration files of the application.

## scripts
This directory contains general scripts.
- **db** : contains db management and migration scripts for the application