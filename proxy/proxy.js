/**
 * Copyright 2022 - Offen Authors <hioffen@posteo.de>
 * SPDX-License-Identifier: Apache-2.0
 */

const api = new Api()

window.addEventListener('message', function handleMessage (evt) {
  (() => {
    if (!evt.data || !evt.data.type || !evt.data.payload) {
      return Promise.reject(new Error('Received no or malformed message, cannot continue.'))
    }
    switch (evt.data.type) {
      case 'QUERY':
        return api
          .get()
          .then((result) => {
            if (!evt.data.payload.scopes.length) {
              return result
            }
            const decisions = evt.data.payload.scopes.reduce((acc, scope) => {
              acc[scope] = null
              if (scope in result.decisions) {
                acc[scope] = result.decisions[scope]
              }
              return acc
            }, {})
            return { decisions }
          })
          .then(wrapResponse('SUCCESS'))
      case 'ACQUIRE':
        return api.get()
          .then(({ decisions: existingDecisions }) => {
            const decisionsToBeTaken = evt.data.payload.scopes.filter((scope) => {
              return !(scope in existingDecisions)
            })
            return requestDecisions(decisionsToBeTaken, function (styles) {
              evt.ports[0].postMessage(wrapResponse('STYLES')(styles))
            })
          })
          .then((decisions) => {
            return api
              .post({ decisions })
          })
          .then(wrapResponse('SUCCESS'))
      case 'REVOKE':
        return api
          .delete()
          .then(wrapResponse('SUCCESS'))
      default:
        return Promise.reject(new Error(`Unsupported message type "${evt.data.type}"`))
    }
  })()
    .catch(wrapResponse('ERROR'))
    .then((response) => {
      if (evt.ports && evt.ports.length > 0) {
        evt.ports[0].postMessage(response)
      }
    })

  function wrapResponse (type) {
    return function (payload) {
      return {
        type,
        payload
      }
    }
  }
})

function Api () {
  this.get = handleResponse(() => {
    return window.fetch(window.location.origin + '/consent', {
      method: 'GET',
      credentials: 'include'
    })
  })

  this.post = handleResponse((body) => {
    return window.fetch(window.location.origin + '/consent', {
      method: 'POST',
      credentials: 'include',
      body: body ? JSON.stringify(body) : undefined
    })
  })

  this.delete = handleResponse((body) => {
    return window.fetch(window.location.origin + '/consent', {
      method: 'DELETE',
      credentials: 'include',
      body: body ? JSON.stringify(body) : undefined
    })
  })

  function handleResponse (fn) {
    return function () {
      return fn.apply(null, [].slice.call(arguments))
        .then((res) => {
          if (res.status === 204) {
            return Promise.resolve(null)
          }
          return res.json()
        })
    }
  }
}

function requestDecisions (scopes, relayStyles) {
  return scopes.reduce((result, scope) => {
    return result.then(decisions => {
      const element = document.querySelector(`[data-scope="${scope}"]`) || document.querySelector('[data-scope="default"]')
      const yes = element.querySelector('[data-yes]')
      const no = element.querySelector('[data-no]')
      return new Promise((resolve, reject) => {
        showElement(element)
        relayStyles({ visible: true })
        setTimeout(() => relayStyles({ rect: element.getBoundingClientRect() }), 0)
        if (!yes || !no) {
          reject(new Error('Could not bind event listeners.'))
          return
        }
        yes.addEventListener('click', handleYes)
        no.addEventListener('click', handleNo)

        function handleYes () {
          unbind()
          resolve(true)
        }
        function handleNo () {
          unbind()
          resolve(false)
        }
        function unbind () {
          no.removeEventListener('click', handleNo)
          yes.removeEventListener('click', handleYes)
        }
      })
        .then((decision) => {
          hideElement(element)
          relayStyles({ visible: false })
          decisions[scope] = decision
          return decisions
        })
        .then(deferBy(0))
    })
  }, Promise.resolve({}))
}

function showElement (el) {
  el.classList.add('show')
}

function hideElement (el) {
  el.classList.remove('show')
}

function deferBy (ms) {
  return function (result) {
    return new Promise((resolve) => {
      setTimeout(() => resolve(result), ms)
    })
  }
}
