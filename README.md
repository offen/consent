<a href="https://www.offen.dev/">
    <img src="https://offen.github.io/press-kit/offen-material/gfx-GitHub-Offen-logo.svg" alt="Offen logo" title="Offen" width="150px"/>
</a>

# consent

__Zero-overhead Consent tooling__

---

Documentation on how to install and run the consent tool is found in the [docs][docs] section.

[docs]: ./docs/README.md

## Development setup

To run the development setup, make sure you have `make`, Docker and `docker-compose` installed.

Clone the repository and start the development server:

```
git clone git@github.com:offen/consent.git
cd consent
make up
```

This gives you a `consent` server running on port 9000 and a test environment that embeds the script running on port 9001.

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
