# Gremlin Images

Failure injectors are based on user-specified images. These images define the behavior of the failure injector. They name the available failure modes and are responsible for applying failures according to the configuration injected via environment variables.

This subproject includes a simple example entrypoint script that provides input validation, and reports the available failure modes if invoked without requisite input. This reporting mechanism is critical for reporting user feedback via Entropy. 

Feedback might include available failure modes, or more specific probability or frequency input requirements.

## An Interface

This document describes the interface contract between Entropy and gremlin images. Any changes to that interface should be documented here and reflected in the entrypoint script.

### Input: Environment Variables

##### ENTROPY_FAILURES

A string name of a failure to inject. The injector image should validate this against a list of named failures that it provides.

##### ENTROPY_FREQUENCY

Typically this field is a positive, base-10 integer representing a duration in seconds. This is symbol is currently used for both the duration of a failure and the frequency with which the dice are rolled. Valid arguments can be made for changing this in the future. 

##### ENTROPY_PROBABILITY

A number, usually a decimal between 1 and zero in base-10 numeric. The value itself will always be a string, so if you want to get fancy with your gremlin implementation you are free to do so. Just make sure that your validation is comprehensive and that the messaging is clear.

##### ENTROPY_RO_ENDPOINT

If the Entropy service exposes an internal read only endpoint with resources for peer discovery or quorum that URL will be provided here. This can be used for failure "election."

### Output: Log Output

The simplest output mechanism is reading STDOUT. The simplest way I can think to test that output for the presence of a system message is to require a special prefix at the beginning of the stream. 

    == Entropy Header Open ==

The first line of STDOUT should contain this line. If the image fails to do so Entropy will assume that the injector has started successfully. However, the system will be unable to provide the user any further insight.

This header block should contain a JSON payload adhering to the JSON spec defined here.

Among other things that payload contains a success code, validation error messages, log format definition and backplane communication line prefix, and a frequency that Entropy should use to check for messages.

The header block should be closed with the following line.

    == Entropy Header Close ==


