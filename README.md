<a href="https://www.offen.dev/">
    <img src="https://offen.github.io/press-kit/offen-material/gfx-GitHub-Offen-logo.svg" alt="Offen logo" title="Offen" width="150px"/>
</a>

# consent

[![CircleCI](https://circleci.com/gh/offen/consent/tree/development.svg?style=svg)](https://circleci.com/gh/offen/consent/tree/development)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Docs](https://img.shields.io/badge/Documentation-docs-blue.svg)][docs]

__Zero-overhead consent tooling__

`consent` aims to be a lightweight solution for managing user consent on websites.
Its most important design goals are:
- no server side persistence of consent decisions
- no need to assign user identifiers or similar, meaning no additional tracking vectors
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

---

Documentation on how to install and run the consent tool is found in the [docs][docs] section.

[docs]: ./docs/README.md

## Quickstart

Deploy the application to a domain like `consent.example.com`.
On the host site `www.example.com`, embed the client script:

```html
<script src="https://consent.example.com/client.js">
```

which exposed `window.ConsentClient`.
In your client side code, construct a new client instance pointing at your deployment and request user consent for the desired scope(s):

```js
const client = new window.ConsentClient({ url: 'https://consent.example.com' })
client
  .acquire('analytics', 'marketing')
  .then((decisions) => {
    if (decisions.analytics) {
      // load analytics data
    }
    if (decisions.marketing) {
      // trigger marketing tools
    }
  })
```

## Building the binary/images yourself

By default, a `linux-x86_64` binary is provided for all releases.
If you need to build a binary for a different platform you can install Go 1.18 or greater, clone the repository and run

```
make
```

which will put the binary in a `bin` directory shortly.

---

Docker images are provided for `amd64`, `arm64` and `arm/v7`, building for other targets is possible via:

```
docker buildx build --platform <your_target> -t offen/offen:<your_tag> .
```

## Development setup

To run the development setup, make sure you have `make`, Docker and `docker-compose` installed.

Clone the repository and start the development server:

```
git clone git@github.com:offen/consent.git
cd consent
make up
```

This gives you a `consent` server running on port 9000 and a test environment that embeds the script running on port 9001.
