# consent

__Zero-overhead Consent tooling__

## Scope and use cases

`consent` aims to be a lightweight solution for managing user consent on websites.
Its most important design goals are:
- no server side persistence of consent decisions
- consent decisions are secured from interference of 3rd party scripts
- users can revoke their consent decisions and any traces at any time by clearing their cookies or using the provided UI
- operators can customize the UI elements in use to match their design

`consent` is a good option in case you want to:
- load 3rd party content like Twitter or Instagram widgets with user consent
- request users' consent for collecting usage statistics
- keep your data footprint as low as possible by never storing any data about consent decisions yourselves

`consent` does not aim to be a drop-in "handle-GDPR-requirements-for-me" solution.
Operators that deploy `consent` are expected to follow regulations that apply themselves.
Ideally, it also acts as a motivation for thinking about what data you really need to collect and which third party services are really required to run your site.

`consent` requires you to be able to configure deploy a simple web server to a dedicated domain.
Linux binaries and a Docker image are provided, or you can build the server for any other platform.
If needed, it can automatically acquire SSL certificates.

`consent` can also be used as a library and be integrated into any web server written in Golang.

## Concepts

TBD.

## Installation and configuration

TBD.

## Embedding the script

TBD.

## Customizing the consent UI

TBD.

## Usage as a library
