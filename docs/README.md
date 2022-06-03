# consent

__Zero-overhead Consent tooling__

<!-- MarkdownTOC -->

- [Concepts and requirements](#concepts-and-requirements)
  - [Designed for client side usage](#designed-for-client-side-usage)
  - [Deploying the server to a sibling domain](#deploying-the-server-to-a-sibling-domain)
  - [Scopes, decisions and consent domains](#scopes-decisions-and-consent-domains)
- [Installation and configuration](#installation-and-configuration)
  - [Running the binary distribution](#running-the-binary-distribution)
  - [Running the Docker image](#running-the-docker-image)
  - [Configuration](#configuration)
    - [`-port` \(`PORT`\)](#-port-port)
    - [`-domain` \(`DOMAIN`\)](#-domain-domain)
    - [`-certs` \(`CERTS`\)](#-certs-certs)
    - [`-ui-copy` \(`UI_COPY`\)](#-ui-copy-ui_copy)
    - [`-ui-button-yes` \(`UI_BUTTON_YES`\)](#-ui-button-yes-ui_button_yes)
    - [`-ui-button-no` \(`UI_BUTTON_NO`\)](#-ui-button-no-ui_button_no)
    - [`-templates-directory` \(`TEMPLATES_DIRECTORY`\)](#-templates-directory-templates_directory)
    - [`-stylesheet` \(`STYLESHEET`\)](#-stylesheet-stylesheet)
- [Embedding the script](#embedding-the-script)
  - [Client usage](#client-usage)
- [Customizing the consent UI](#customizing-the-consent-ui)
  - [Customizing the default UI](#customizing-the-default-ui)
  - [Providing custom content for scopes](#providing-custom-content-for-scopes)
  - [Positioning the consent UI](#positioning-the-consent-ui)
- [Usage as a library](#usage-as-a-library)

<!-- /MarkdownTOC -->


## Concepts and requirements

### Designed for client side usage

`consent` is designed to be used in the context of client side code.
Use cases like deferring loading of third party content are well supported, where usage of consent decisions on the server is an explicit non-goal.

### Deploying the server to a sibling domain

`consent` is using 1st Party Cookies to store user's consent decisions.
To enable this mechanism, you need to deploy the `consent` server to a sibling domain, i.e. if you plan to use the tool on `www.example.com`, `consent` should be served on a domain like `consent.example.com`.
The tool can serve any number of domains at once, so it's possible to use the same deployment for multiple domains at once.

The tool expects to be served via HTTPS so that it can use the `Secure` attribute on the cookies issued.
In case you cannot provide a certificate yourself, the server can acquire a free and automatically renewed certificate from Let's Encrypt in case you provide the matching configuration.

### Scopes, decisions and consent domains

In the context of `consent`, a __decision__ is the answer to a consent request.
Such a decision can only ever be __yes__(`true`) or __no__(`false`).
Each decision belongs to a __scope__, which is a simple string identifier, e.g. `twitter`, `marketing`, `analytics` or anything else that you might want to request user consent for.
Additionally, consent decisions are bound to the referring domain, i.e. a consent decision that has been taken on `www.example.com` is not valid for `app.example.com`.
Each domain needs to request consent again, preventing accidental or inadvertent user consent on such domains.

## Installation and configuration

`consent` runs a simple web server that is expected to be bound to a port and exposed to the internet so it can answer HTTP(S) requests.

### Running the binary distribution

`consent` is built into a single binary that contains all required assets.
Head to the [Releases][releases] page to download the Linux binary or build it yourself.

An example `systemd` service file could look like:

```
[Unit]
Description=Offen Consent Tool

[Service]
ExecStart=/usr/local/bin/consent
Restart=always

[Install]
WantedBy=multi-user.target
```

[releases]: https://github.com/offen/consent/releases

### Running the Docker image

All releases of `consent` are also published to Docker Hub as `offen/consent`.
To run the server in Docker, use:

```
docker run --rm -p 80:80 -p 443:443 offen/consent:latest
```

### Configuration

Configuring the `consent` server can happen in three different means, using the following precendence:
1. command line flags
1. environment variables
1. a configuration file in YAML format, specified by passing `-config`

The following options are available:

#### `-port` (`PORT`)

The port the server binds to.
Defaults to 8000 (80 when run in Docker).
This has no effect in case `domain` is set.

#### `-domain` (`DOMAIN`)

The domain to acquire and manage SSL certificates for.
This is done using LetsEncrypt and HTTP challenges, so make sure to allow Traffic on both port 80 and 443 if you are planning to use this feature.
The default value is the empty string, which makes the server listen to the default port.

#### `-certs` (`CERTS`)

The file system location to use for caching SSL certificates.
Defaults to `/var/www/.cache`.
When using this, make sure the location is readable and writable for the user owning the service process.

#### `-ttl` (`TTL`)

The duration for which a cookie is valid after it has been issued.
Defaults to `4464h`.
When using this, pass a string that can be parsed by Golang's [`time.ParseDuration`][duration].

[duration]: https://pkg.go.dev/time#ParseDuration

#### `-ui-copy` (`UI_COPY`)

In case you are using the `default` scope provided by the application, this option sets the copy used in the consent UI's description paragraph.
No default.

#### `-ui-button-yes` (`UI_BUTTON_YES`)

In case you are using the `default` scope provided by the application, this option sets the copy used in the consent UI's "Yes" button.
No default.

#### `-ui-button-no` (`UI_BUTTON_NO`)

In case you are using the `default` scope provided by the application, this option sets the copy used in the consent UI's "No" button.
No default.

#### `-templates-directory` (`TEMPLATES_DIRECTORY`)

This option allows you to provide custom HTML templates for the scopes you are planning to request consent for.
Each scope is expected to be put in a `<scope identifier>.html` file and will be inlined into the host document.

#### `-stylesheet` (`STYLESHEET`)

This option allows you to pass the location of CSS file that is inlined into the host document.
You can use this option to apply styling to the scope elements in `templates-directory`

## Embedding the script

### Client usage

## Customizing the consent UI

### Customizing the default UI

### Providing custom content for scopes

### Positioning the consent UI

## Usage as a library

`consent` can also be used as a Go library.
Documentation is available in [godoc][] format.

[godoc]: https://pkg.go.dev/github.com/offen/consent
