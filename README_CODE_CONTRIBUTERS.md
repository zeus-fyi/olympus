## Zeus Technical Commandments

Zeus has issued the following proclamations

* Above all, the engineering UX should be considered first, and the consumer UX second.
  * Code should be modular, extensible, readable, and have the right amount of testing in them.
  * If your code is hard to use and understand, it isn't far from just being completely useless.
  * If your code doesn't make it very easy for the engineers and outside consumers, it's not good enough.
  * Architecture Principles:
    * When architecting, or writing code, think about and consider how maintenance and testing will work
      * Optimize UX, Maintaining
        * Easy to maintain, and ideally should be able to self-maintain any non-novel issues with automation.
        * Use automation liberally via Iris to automate fixing, upgrading, troubleshooting.
      * Optimize being the base of the ecosystem pyramid when it's a core module component
        * Meaning, always thing about how/what might be needed for others, and be flexible for unexpected needs
        * Make it easy to build into/out of your designs
        * Make it API first
      * Make things stateless when you can, and for even better, use serverless functions whenever possible.
      * Optimize tradeoffs
        * Consider bottlenecks and costs.
          * Eg. Cost CPU > RAM > SSD > HDD. 
          * Latency. Resilient to failure. 
          * Use performance data from the system to allow dynamic automated infrastructure changes to happen underneath. 
  * Coding Principles:
    * Avoid creating code file pages longer than 1 page when possible, be liberal with more pages and organizations.
    * Functions should strive to be entirely visible without scrolling.
    * Use naming pre-fixing to take advantage of filename sorting
        * eg config-map-part1, config-map-part2
    * Functions should only create one thing.
        * Eg function that calls multiple other functions has created one thing. Use another function to mutate things.
    * If something is expected to take a long time to execute (eg a complex integration test)
      * It should strive to not block the user/engineer and should run async -> notify when done
    * Code should execute very fast when synchronous. So you aren't waiting around all day on this.
        * Unless there's a very good reason, like some external blockchain interaction
        * You should probably decompose your code if it's too slow
    * Test suites for a directory should execute within 5 min or less, and be parallel when possible
        * Modules be fully testable (or at least 95%) within 15 minutes.
        * Integration testing can take longer, just same concepts though.
    * If a Sr. Engineer can't figure out how to safely test, add, modify, delete a module within 3 days or less
        * Add a README to explain things further
* Applications & modules should NOT be coupled to others without good reason
* API first designs, everything should be programmable
* Strive to write code using meta programming, or templating
    * Meaning, if you can use a codegen toolkit use this because this allows you to
        * Refactor code instantly
        * Write code 1000x faster that is reliable, tested, pre-styled, and optimized
        * Write code on the fly, e.g. an application that can write code
        * Autogenerate open api files, and api documentation, etc
    * If the tooling to efficiently do this isn't available/mature enough, these tips will greatly assist you
        * Write code in a style that matches what you'd expect a computer to write.
        * Use common styling and naming conventions
        * Limit each code file scope as much as possible
        * Come up with uniform templated test_suite/func/struct to copy/paste while prototyping
            * Use this style to codify using the codegen tools
* Codify complex steps using Iris (Temporal/Confluent Engines)
    * Eg build -> test -> message
* Make everything cloud provider/vendor independent whenever possible, maximize portability.
* Be efficient with infrastructure
  * Questions to ask yourself
    * Does this need this class of performance? 
      * Example: I put this data in cold storage which is much, much cheaper since it's rarely accessed
    * Does another cloud provider offer an equivalent for a lot less?
    * Is this being used? 
    * Is this scalable, do costs go down per unit with scale by a lot? ideally it should be slightly linear
* Don't use helm charts for mature internal apps, use Zeus
* Don't use or add any terrible languages like bash, jsonnet, php, etc.
    * Unless you can autogenerate it entirely from a structured language
* Don't write ANY cli tooling, or Zeus WILL banish you from Olympus (exceptions rarely granted)
    * Cli tooling is dumb, don't be a dummy. Why is it dumb? Glad you asked, here's why.
        * Cli tooling implies you're going to be asking a mere mortal to do engineering by hand
            * Doing anything by hand means you can easily mistype something, and leads to arthritis
            * Adding more Humans-In-The-Loop is a terrible bottleneck to efficiency
        * Cli tooling is difficult to manage version control for managing breaking changes
        * Cli tooling is annoying and makes it more difficult to write good tests
        * Cli tooling is hard to chain and build more advanced tooling from
            * It also has a side effect of encouraging shell script creation (don't make Zeus even more mad now)
* Don't write ANY yaml by hand for anything, unless these exceptions apply or Zeus will hit you with a thunderbolt
    * Updating or maintaining an externally sourced helm chart used in production
    * Prototyping internal helm chart (dev/staging only) prior to insertion into Zeus
    * Simple local configs like docker compose files
* Don't write shell scripts or Zeus will turn you to stone.
    * Shell scripts and bash languages have horrific UX and cause errors that are hard to solve.
    * See cli tooling banishment reasoning for more reasons.
    * Allowable exceptions:
        * Using an externally created shell script
        * Internal exceptions
            * You MUST document where/what/why so this can be removed later on when a better solution exists
            * It contains no conditional logic (eg if, else, etc), and <= 10 lines of code.
            * Only and only if for some reason it is at least 10x more efficient way than other approaches.
            * Any other exceptions are unlikely to be granted 

### Useful Development Tips

These files are .gitignored by default here, and meant for creating local notes

* ****TODO.txt****
    * Use to store notes about dev plans
* ****SCRATCH_PAD.txt****
    * Use to code, links, etc to quickly reference or store copy/paste examples

##### Environment Specific Interactions Setups

Wrappers for env specific interactions from a local user should be done in `/pkg/utils/env`
The wrappers should link to the env specific configs, like connecting to a staging postgres instance.

* Local
* Development
* Staging
* Production
    * Must access sessions with user roles if possible
        * User role permission should be set to limit ability to limit scope and prevent destructive actions like table truncations
    * Should only use read-only interactions by default
    * In limited cases where it makes sense, such as inserting data into the production database, it should use a specially created role that does not allow deletions or table truncations

##### Local Testing Resource Setups

Find useful makefile commands to start services in docker as needed

* Postgres
* Redis
* Temporal

#### Olympus Language Guidelines & Conventions

* Use these languages, unless a strong case for using an alternative is made like no SDK language
    * Common Backend (version preference)
        * Golang (go 19.2)
        * Python (use latest)
        * Typescript (use nvm)
            * For when you are stuck with writing JS.
    * Frontend
        * TBD. Possibly using auto-generated go-templates that get controlled by JS React module wrappers
    * Data languages
        * Write the SQL for everything, don't use ORMs ever.
        * Codify handwritten SQL pattern abstractions to then use autogen
            * Ideally autogen path does this from a common source
                * Creates tables, migrates data
                * Generates Go structs from table schemas with db tags and json tags
                * Generates common SQL interactions with test suites in Golang (non-exhaustive examples below)
                    * Read/Insert/Update/Select Data
                    * Simple Join Queries

### Code Directory Guide ###

Monorepo that is modular by design

### Datastores

Technology Stack:

* Postgres 13.4 (due to temporal dependency)
* Redis (latest)

Common Directory Structure Convention

* datastore_name (eg. Postgres)
    * datastore_apps (to store app interfacing logic, structs, app)
        * COMMON_FILE: datastore_class_name (common base file to connect & interact with datastore)
    * datastore_docker (to store local docker setups, and migration files when applicable)
        * COMMON_FILE: docker-compose-datastore.yml
