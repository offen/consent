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
    this.send = this.injectIframe(url)
    this.targetOrigin = new window.URL(url).origin
  }

  injectIframe (url) {
    var proxy = document.createElement('iframe')
    proxy.src = url

    proxy.style.display = 'none'
    proxy.setAttribute('frameBorder', '0')
    proxy.setAttribute('scrolling', 'no')
    proxy.setAttribute('title', 'Consent Proxy')

    var elementId = 'consent-proxy-' + Math.random().toString(36).slice(2)
    proxy.setAttribute('id', elementId)

    var iframe = new Promise(function (resolve, reject) {
      proxy.addEventListener('load', function (e) {
        function postMessage (message) {
          return new Promise(function (resolve, reject) {
            var origin = new window.URL(proxy.src).origin
            message.host = message.host || '#' + elementId

            var messageChannel = new window.MessageChannel()
            messageChannel.port1.onmessage = function (event) {
              var responseMessage = event.data || {}
              switch (responseMessage.type) {
                case 'ERROR':
                  var err = new Error(responseMessage.payload.error)
                  err.originalStack = responseMessage.payload.stack
                  err.status = responseMessage.payload.status
                  reject(err)
                  break
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
    return this.send.then(send => {
      return send({ type, payload: { scopes } })
    })
  }
}

module.exports = Client
