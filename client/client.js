/**
 * Copyright 2022 - Offen Authors <hioffen@posteo.de>
 * SPDX-License-Identifier: Apache-2.0
 */

class Client {
  constructor (options = {}) {
    this.proxy = new EmbeddedProxy(options.url, options.host, options.ui)
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
  constructor (url, host, ui) {
    this._send = this.injectIframe(url, host, ui)
    this.targetOrigin = new window.URL(url).origin
  }

  injectIframe (
    url,
    host = document.body,
    options = {}
  ) {
    options = Object.assign({
      className: 'consent',
      styles: { margin: 'auto', position: 'fixed', bottom: '1em', left: '0', right: '0' }
    }, options)
    const proxy = document.createElement('iframe')
    proxy.src = url + '/proxy'

    proxy.style.display = 'none'
    proxy.setAttribute('frameBorder', '0')
    proxy.setAttribute('scrolling', 'no')
    proxy.setAttribute('title', 'Consent Proxy')

    const elementId = 'consent-proxy-' + Math.random().toString(36).slice(2)
    proxy.setAttribute('id', elementId)
    proxy.classList.add(options.className)
    for (const prop in options.styles) {
      proxy.style[prop] = options.styles[prop]
    }

    const iframe = new Promise(function (resolve, reject) {
      proxy.addEventListener('load', function (e) {
        function postMessage (message) {
          return new Promise(function (resolve, reject) {
            const origin = new window.URL(proxy.src).origin
            message.host = message.host || '#' + elementId

            const messageChannel = new window.MessageChannel()
            messageChannel.port1.onmessage = function (evt) {
              const responseMessage = evt.data || {}
              switch (responseMessage.type) {
                case 'STYLES':
                  if ('visible' in evt.data.payload) {
                    proxy.style.display = evt.data.payload.visible
                      ? 'block'
                      : 'none'
                  }
                  if ('rect' in evt.data.payload) {
                    proxy.setAttribute('width', evt.data.payload.rect.width)
                    proxy.setAttribute('height', evt.data.payload.rect.height)
                  }
                  break
                case 'ERROR': {
                  const err = new Error(responseMessage.payload.message)
                  err.originalStack = responseMessage.payload.stack
                  err.status = responseMessage.payload.status
                  reject(err)
                  break
                }
                default:
                  resolve(responseMessage.payload)
              }
            }
            messageChannel.port1.onmessageerror = function (err) {
              reject(err)
            }
            proxy.contentWindow.postMessage(message, origin, [messageChannel.port2])
          })
        }
        resolve(postMessage)
      })
      proxy.addEventListener('error', function (err) {
        reject(err)
      })
    })

    switch (document.readyState) {
      case 'complete':
      case 'loaded':
      case 'interactive':
        host.appendChild(proxy)
        break
      default:
        document.addEventListener('DOMContentLoaded', function () {
          host.appendChild(proxy)
        })
    }
    return iframe
  }

  send (type, ...scopes) {
    return this._send.then(send => {
      return send({ type, payload: { scopes } })
    })
  }
}

window.ConsentClient = Client
