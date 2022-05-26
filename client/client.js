/**
 * Copyright 2022 - Offen Authors <hioffen@posteo.de>
 * SPDX-License-Identifier: Apache-2.0
 */

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
    this._send = this.injectIframe(url)
    this.targetOrigin = new window.URL(url).origin
  }

  injectIframe (url) {
    const proxy = document.createElement('iframe')
    proxy.src = url + '/proxy'

    proxy.style.display = 'none'
    proxy.setAttribute('frameBorder', '0')
    proxy.setAttribute('scrolling', 'no')
    proxy.setAttribute('title', 'Consent Proxy')

    const elementId = 'consent-proxy-' + Math.random().toString(36).slice(2)
    proxy.setAttribute('id', elementId)

    const iframe = new Promise(function (resolve, reject) {
      proxy.addEventListener('load', function (e) {
        function postMessage (message) {
          return new Promise(function (resolve, reject) {
            const origin = new window.URL(proxy.src).origin
            message.host = message.host || '#' + elementId

            const messageChannel = new window.MessageChannel()
            messageChannel.port1.onmessage = function (event) {
              const responseMessage = event.data || {}
              switch (responseMessage.type) {
                case 'ERROR': {
                  const err = new Error(responseMessage.payload.error)
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
        document.body.appendChild(proxy)
        break
      default:
        document.addEventListener('DOMContentLoaded', function () {
          document.body.appendChild(proxy)
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

window.Client = Client
