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
      return api.get(evt.data.payload).then(wrapResponse('SUCCESS'))
    case 'ACQUIRE':
      return api.post(evt.data.payload).then(wrapResponse('SUCCESS'))
    case 'REVOKE':
      return api.delete(evt.data.payload).then(wrapResponse('SUCCESS'))
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
})

function Api () {
  this.get = ()  => {
    return window.fetch(window.location.origin + '/consent', {
      method: 'GET',
      credentials: 'include'
    })
      .then(handleResponse)
  }

  this.post = (body)  => {
    return window.fetch(window.location.origin + '/consent', {
      method: 'POST',
      credentials: 'include',
      body: body ? JSON.stringify(body) : undefined
    })
      .then(handleResponse)

l }

  this.delete = (body)  => {
    return window.fetch(window.location.origin + '/consent', {
      method: 'DELETE',
      credentials: 'include',
      body: body ? JSON.stringify(body) : undefined
    })
      .then(handleResponse)
  }

  function handleResponse (res) {
    if (res.status === 204) {
      return Promise.resolve(null)
    }
    return res.json()
  }
}

function wrapResponse (type, payload) {
  if (!payload) {
    return function (payload) {
      return wrapResponse(type, payload)
    }
  }
  return {
    type: type,
    payload: payload
  }
}
