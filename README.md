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

This gives you a `consent` server running on port 9000 and a test environement that embeds the script running on port 9001.
