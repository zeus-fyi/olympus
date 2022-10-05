# README #

Olympus is built from many modules that work in cohesion as deployable microservices.

The default philosophy to use is that applications/modules should not be coupled to others

### How do I get set up? ###

* Summary of set up
* Configuration
* Dependencies
* Database configuration
* How to run tests
* Deployment instructions

### Interactions with Environments

Wrappers for env specific interactions should be done in `/pkg/utils/env`

* Local
* Development
* Staging
* Production
  * Should only use read-only interactions by default
  * In limited cases where it makes sense, such as inserting data into the production database, it should use a specially created role that does not allow deletions or table truncations


##### Local Testing Resource Setups

Find useful makefile commands to start services in docker as needed

* Postgres
* Redis
* Temporal

### Code Directory Guide ###

Monorepo that is modular by design

### Apps

Where full applications are built

  * Apollo
      * Observes and captures blockchain data
  * Zeus
      * Dynamically controls infrastructure and applications

#### Configs

Where all configs should be placed by default. Uses gitignore to prevent commit of sensitive values

#### Pkg

Libraries to build applications from

* Codegen
    * For generating code and creating templates
* Zeus
    * Where Kubernetes API commands are built

#### Devops

Where flux configurations are kept

#### Helm

Stores helm chart definitions. Flux uses this git repo to track and deploy helm chart changes. Chart museum is
mostly used for downloading external charts.

#### Sandbox

Stores data for mocking, experimenting, prototyping
