module.exports = function Api () {
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
