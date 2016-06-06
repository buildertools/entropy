# Entropy 

Entropy is a failure injectior orchestration service for Docker platforms. 

## State of the Project

The project is currently PoC quality with no tests and a crude software design. Most existing code will be replaced before a proper beta.

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
