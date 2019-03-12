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
                console.log('Status Body parse error: ', url, '\n', response);
                return false
            }
        } else {
            return response
        }
    }
    console.log('Status NotOk: ', url, '\n', response);
    return false;
}

const _get = async (url, parseResults) => _callApi(url, {method: 'get'}, parseResults);

const callLoginApi = async (email, password) => _callApi('login', {
    method: 'post',
    body: JSON.stringify({
        Email: email,
        Password: password
    })
})

const callAuthcheckApi = async () => _get('authcheck')
const callUpgradeApi = async () => _get('upgrade', true)
const callUpgradeCheckApi = async () => _get('upgradecheck', true)
const callDashboardApi = async () => _get('dashboard', true)
const callLogoutApi = async () => _get('logout')

const callApi = async () => {
    // DEPRECATED - test util
    const response = await fetch('/api');
    if (response.status !== 200) throw Error(response);
    return response;
};

export {
    callDashboardApi,
    callApi,
    callAuthcheckApi,
    callLoginApi,
    callUpgradeApi,
    callUpgradeCheckApi,
    callLogoutApi
}
