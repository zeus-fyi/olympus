# Zeus Models #

High Level Explainer

### What is this repository for? ###

* Zeus Model Definitions and Interfaces

### Technology Stack ###

* Postgres 13.4 (Temporal Dependency)

### Table of Contents ###
* Database

## Database

### Models

* Models

Model structs have two classes
* autogen
* manual

Use auto generation for creating base model structs using the codegen makefile

makefiles/codegen/Makefile

The manual structs should be higher level wrapper structs over the autogen ones
