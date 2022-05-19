/**
 * Copyright 2022 - Offen Authors <hioffen@posteo.de>
 * SPDX-License-Identifier: Apache-2.0
 */

var api = new Api()

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
          var decisions = evt.data.payload.scopes.reduce((acc, scope) => {
            acc[scope] = null
            if (scope in result.decisions) {
              acc[scope] = result.decisions[scope]
            }
            return acc
          }, {})
          return { decisions: decisions }
        })
        .then(wrapResponse('SUCCESS'))
    case 'ACQUIRE':
      return requestDecisions(evt.data.payload.scopes)
        .then((decisions) => {
          return api
            .post({ decisions: decisions })
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
        type: type,
        payload: payload
      }
    }
  }
})

function Api () {
  this.get = handleResponse(()  => {
    return window.fetch(window.location.origin + '/consent', {
      method: 'GET',
      credentials: 'include'
    })
  })

  this.post = handleResponse((body)  => {
    return window.fetch(window.location.origin + '/consent', {
      method: 'POST',
      credentials: 'include',
      body: body ? JSON.stringify(body) : undefined
    })
l })

  this.delete = handleResponse((body)  => {
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

function requestDecisions (scopes) {
  var decisions = scopes.reduce((acc, next) => {
    acc[next] = true
    return acc
  }, {})
  return Promise.resolve(decisions)
}
