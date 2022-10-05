# Olympus #

Olympus is built from many modules that work in cohesion as deployable microservices.

### Cloud Providers

Digital Ocean
  * Kubernetes
  * Temporal
  * Blob Storage
  * Docker Image Registry
Azure
  * Postgres
  * Confluent

### High Level Infrastructure Stack

Infrastructure as Code
  * Kubernetes
    * Sources
      * Olympus specific workloads for dynamic k8
      * Helm
  * Flux

Datastores
  * Postgres
  * Redis

Orchestration
  * Temporal

Messaging/Streaming
  * Confluent

### Configs

Where all configs should be placed by default. Uses gitignore to prevent commit of sensitive values

### Apps

Where full applications are built

* Apollo
  * Observes and captures blockchain data
* Zeus
  * Dynamically controls infrastructure and applications

### Pkg

Libraries to build applications from. Each package component may have a README.md with more verbose information.

* Zeus
    * Where Kubernetes API commands are built and other dynamical api controllers
* Utils
  * Library for common code sharing. When versioning is needed follow this guideline
    * Create a wrapper struct over the module component (contrived example for adding a v1 logging feature)
    * Structs should use a common versioned log to track when it's being used and what dependency paths exist
      * Base Struct Logging Fields to Inherit
        * V0_Logging_Library
        * V1_Logging_Library
      * Create Versioned Base Struct to Inherit From Parent(s)
        * Now create V1 functions like: func (l LoggerV1) () {}
    * Why do this?
      * Using this style you'll never have to update things by hand, it'll be implied
      * You won't break anything
  * Should structure library like this when it becomes more complex (top-bottom)
    * Complex
      * Peak - relies on multiple middle level library components or an external dependency
      * Middle - built from many base library components
    * Base
      * Should strive to be stand-alone building blocks and be highly stable
      * Changes should create a new version reference if they're referenced by any production service in these cases
        * Not backwards compatible or any change in expected behavior
* Codegen
  * For generating code and creating templates
  * Go code-gen technology stack
    * Jennifer: https://github.com/dave/jennifer
    * ToJen: https://github.com/aloder/tojen
      * For parsing a file into a Jennifer config
      * Makes prototyping new code gen abstractions a lot easier
          
### Devops

Where flux configurations are kept

### Helm

Stores helm chart definitions. Flux uses this git repo to track and deploy helm chart changes. Chart museum is
mostly used for downloading external charts.

### Docker

Stores docker build configurations. In the future, this should all be called by Circe to build, then test via Ares
using orchestrations, and then storing results via Hestia (persistent data sources)

CircleCI is integrated into GitHub, but it'll probably be deprecated once Circe comes online

### Sandbox

Stores data for mocking, experimenting, prototyping
