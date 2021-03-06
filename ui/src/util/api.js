/**
 * Api Helpers
 */
const callLoginApi = async (email, password) => _callApi('login', {
    method: 'post',
    body: JSON.stringify({
        Email: email,
        Password: password
    })
})

const callLogoutApi = async () => _get('logout')

const callAuthcheckApi = async () => _get('authcheck')

const callUpgradeApi = async () => _get('upgrade', true)

const callDashboardApi = async () => _get('dashboard', true)

// Internals
/**
 * Returns a call with auth creds to api/url
 * @param {string} url url endpoint to proxy to, will prepend /api/
 * @param {*} extensions any extension options to merge into our req
 * @param {*} parseResults Should we parse results? if false will not attempt to read response json stream
 */
const _callApi = async (url, extensions, parseResults=false) => {
    const response = await fetch(`/api/${url}`, Object.assign({
        headers: {
            Accept: 'application/json',
            'Content-type': 'application/json'
        },
        credentials: 'same-origin'
    }, extensions));

    if (response.status === 200) {
        if (parseResults) {
            try {
                let finalResults = await response.json()
                return finalResults
            } catch (e) {
                console.log('Status Body parse error: ', url, '\n', response)
                return false
            }
        } else {
            return response
        }
    }
    console.log('Status NotOk: ', response.status, url, '\n', response)
    return false
}

const _get = async (url, parseResults) => _callApi(url, {method: 'get'}, parseResults);

export {
    callDashboardApi,
    callAuthcheckApi,
    callLoginApi,
    callUpgradeApi,
    callLogoutApi
}
