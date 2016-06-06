# Entropy 

Entropy is a failure orchestration microservice for Docker platforms. 

Entropy allows a user to define failure injection policies for containers or services running at a target Docker API endpoint. A failure injection policy describes a failure type, frequency/duration, probability, and container selector. Each policy is applied to all current and future containers matching the selector. Policies are manifest by the creation, scheduling, and lifecycle management of sidekick containers called injectors.

### Failures

Entropy failures are pluggable, but by default there are four network failures available:

* Latency
* Network Partition
* Packet Reordering
* Packet Loss

### Selectors

Policies are applied to all current and future containers at the target endpoint that match the provided selector. Currently selectors are limitied to a single label specification (as in Docker container labels). These may take the form: K, or K=V.

This is sufficient to create broad, cross-cutting, or specific failure injection policies for systems with healthy taxonomies.

## State of the Project

The project is currently PoC quality with no tests and a crude software design. Most existing code will be replaced before a proper beta. Documentation is currently very light as so much will change before the next release. If you'd like to see how to use it, checkout the demo.

## Licencing

The project is currently unlicenced. For the time being, Jeff Nickoloff and All in Geek Consulting Services, LLC retain all copyright to this material. This will change in the next few releases.

## Demo

Checkout [buildertools/entropy-demo] for a demo guide.

## Contributing

The service source is Go, and the scaffolding is a mix of Docker tooling, Go, JavaScript, and Makefile. The whole environment is designed to be approachable and hackable for those familiar with Docker.

A developer will almost always want a local build of the service running while they are developing. Rather than leave setup to contirbutors this project uses a simple Docker based development environment and iteration scaffold. The following set of commands will handle almost every development phase. This system automates the whole lifecycle on source change events. It also manages a fully tooled local runtime environment including time series database, dashboarding, and key-value stores.

#### Get the Code and Get Started

    # Get the code
    git clone https://github.com/buildertools/entropy.git
    
    # Prepare the first time to get tooling
    cd entropy
    make prepare

#### Local Iteration

    # Start iterating
    make iterate
    
    # Do your work:
    #  Use the running service:
    #    open http://localhost:7890
    
    #  Watch the build logs:
    docker-compose logs entropy
    
    #  Watch the functional / integration tests:
    docker-compose logs integ
    
    #  Watch the service metrics:
    #    open http://localhost:3000
    
    # Stop working for the day:
    make stop

#### Updating Dependencies or Build System

    make stop
    # Make your changes
    make iterate

#### Ship it

    # Produce a versioned release image
    make release

Hook developers should add service definitions for target endpoints to the development environment.
