<a href="https://www.offen.dev/">
    <img src="https://offen.github.io/press-kit/offen-material/gfx-GitHub-Offen-logo.svg" alt="Offen logo" title="Offen" width="150px"/>
</a>

# `consent` manual

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
    - [`-ttl` \(`TTL`\)](#-ttl-ttl)
    - [`-ui-copy` \(`UI_COPY`\)](#-ui-copy-ui_copy)
    - [`-ui-button-yes` \(`UI_BUTTON_YES`\)](#-ui-button-yes-ui_button_yes)
    - [`-ui-button-no` \(`UI_BUTTON_NO`\)](#-ui-button-no-ui_button_no)
    - [`-templates-directory` \(`TEMPLATES_DIRECTORY`\)](#-templates-directory-templates_directory)
    - [`-stylesheet` \(`STYLESHEET`\)](#-stylesheet-stylesheet)
- [Embedding the script](#embedding-the-script)
    - [`options.origin`](#optionsorigin)
    - [`options.host`](#optionshost)
    - [`options.ui`](#optionsui)
  - [Client usage](#client-usage)
    - [`client.acquire(...scopes)`](#clientacquirescopes)
    - [`client.query(...scopes)`](#clientqueryscopes)
    - [`client.revoke(...scopes)`](#clientrevokescopes)
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
Such a decision can only ever be __yes__ (`true`) or __no__ (`false`).
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

Configuring the `consent` server can happen in three different means, using the following precedence:
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

The application serves a script file at `/client.js` that needs to be embedded on any site you want to use the tool on.

```html
<script src="https://consent.example.com/client.js"></script>
```

This script provides a global `ConsentClient` class you can instantiate in your code:

```js
const client = new window.ConsentClient(/* options */)
````

In case the `ConsentClient` identifier was already defined, you can restore its old value using the `noConflict` method:

```js
const NoConflictConsentClient = window.ConsentClient.noConflict()
const client = new NoConflictConsentClient(/* options */)
```

The constructor receives an optional options object with the following properties:

#### `options.origin`

The origin (e.g. `https://consent.example.com`) of the target server.
If no value is given, the origin of the script is used.

#### `options.host`

A DOM element that will be used as the parent element of the consent UI.
Child elements of the host will be removed when the iframe is injected.
This can be used to add fallback content for users that have JavaScript disabled or similar.
Defaults to `document.body`.

#### `options.ui`

`ui` is an object defining:

`style` properties that define how the consent UI is styled and positioned.
Any valid style property can be used.
Default values are:
```js
{
  style: {
    margin: 'auto',
    position: 'fixed',
    'max-width': '479px',
    bottom: '1em',
    left: '0',
    right: '0'
  }
}
```

A `className` value, defining a class name that is applied to the `iframe` element containing the consent UI.
This defaults to `consent-iframe`.

```js
{ className: 'consent-iframe' }
```

### Client usage

A client instance exposes three methods.
All methods receive an arbitrary number of decision `scopes`.
In case no arguments are given, the call assumes the operation is applicable to all scopes.
All methods return a Promise resolving with an object containing the user's consent decisions.

#### `client.acquire(...scopes)`

Call `acquire` in case you want to acquire consent for a certain scope.
In case the user already has previously made a decision for a scope, its known state is returned.
In case the decision is still pending, the consent UI will be shown and the returned Promise resolves once a decision has been taken.

#### `client.query(...scopes)`

Call `query` to query for the existing decision on a certain scope.

#### `client.revoke(...scopes)`

Call `revoke` to revoke the existing decision on a certain scope.
This resets the consent decisions to a state matching the one of a new user.

## Customizing the consent UI

`consent` isolates the consent UI into an `iframe` element in order to shield it from unwanted access by third party scripts running in the context of your website. UI can be customized in two different ways:
- provide the text content to be used in the default UI
- provide custom HTML and CSS content to be used for rendering the consent UI

### Customizing the default UI

If you do not need custom styling, you can use the default banner provided.
Either by setting env vars or by providing values to the CLI flags, set the following:
- `UI_COPY` contains the text used to explain what consent is acquired for
- `UI_BUTTON_YES` contains the text displayed on the button for giving consent
- `UI_BUTTON_NO` contains the text displayed on the button for declining consent

### Providing custom content for scopes

`consent` also allows you to provide entirely custom content for your UI.
To provide custom content, provide a directory containing one `.html` file for each scope you want to define using `TEMPLATES_DIRECTORY`.
Such a file is expected to contain an HTML snippet that is then wrapped in a `.consent-scope` element.
Each of these snippets are expected to contain __exactly one__ clicakble element for giving consent and __exactly one__ clickable element that can be used for declining.
These elements __must__ be marked up using `data-yes` and `data-no` attributes.

An example scope would look like:

```
<p>Do you allow collection of anonymous usage data while visiting this site?</p>
<button data-yes>Yes</button>
<button data-no>No</button>
```

A custom stylesheet for styling these scopes can be provided using the `STYLESHEET` option.

In case the client requests a decision for a scope that has not been defined by you, the `default` scope will be used.
You can override the predefined `default` scope by providing a `default.html` file in the `TEMPLATES_DIRECTORY`.

### Positioning the consent UI

Positioning the `iframe` containing the consent UI on the host page is done using the client script.
You can:
1. pass an arbitrary DOM element as `options.host` that acts as the `iframe`'s parent element
2. pass an object containing style attributes that will be applied to the `iframe` element using `options.ui.styles`

In case this does not meet your requirements, you can also position and style the `iframe` element in your own code.
To do so, pass `null` to `options.ui.styles` and target the `iframe` element using the `.consent-iframe` (or whatever you define yourself in `options.ui.className`) selector.

## Usage as a library

`consent` can also be used as a Go library.
Documentation is available in [godoc][] format.

[godoc]: https://pkg.go.dev/github.com/offen/consent
