class Client {
  constructor (options = {}) {
    this.proxy = new EmbeddedProxy(options.url)
  }

  acquire (...scopes) {
    return this.proxy.send('ACQUIRE', ...scopes)
  }

  query (...scopes) {
    return this.proxy.send('QUERY', ...scopes)
  }

  revoke (...scopes) {
    return this.proxy.send('REVOKE', ...scopes)
  }
}

class EmbeddedProxy {
  constructor (url) {
    this.iframe = this.injectIframe(url)

    this.targetOrigin = new window.URL(url).origin
  }

  injectIframe (url) {
    return Promise.resolve(window)
  }

  send (type, ...scopes) {
    return this.iframe.then(el => {
      el.postMessage({ type, scopes }, this.targetOrigin)
    })
  }
}

module.exports = Client
