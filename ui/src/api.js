
import React from 'react';
import { Redirect } from "react-router-dom";

class RouterHack extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            redirectToReferrer: false,
            loginSuccess: false
        }
    }
    shouldComponentUpdate(nextProps) {
        if (nextProps.loginSuccess !== this.state.loginSuccess ||
            nextProps.redirectToReferrer !== this.state.redirectToReferrer) {

            this.setState({
                redirectToReferrer: nextProps.redirectToReferrer,
                loginSuccess: nextProps.loginSuccess
            });
            return true;
        }
        return false;
            
    }
    render() {
        let { redirectToReferrer, loginSuccess } = this.state;

        if (redirectToReferrer) {
            this.setState({redirectToReferrer: false});
            return <Redirect to='/login' />;
        } else if (loginSuccess) {
            this.setState({loginSuccess: false});
            return <Redirect to='/dashboard' push />;
        } else {
            return <div></div>
        }
    }
}

const _callApi = async (url, extensions) => {
    const response = await fetch(`/api/${url}`, Object.assign({
        headers: {
            Accept: 'application/json',
            'Content-type': 'application/json'
        },
        credentials: 'same-origin'
    }, extensions));

    if (response.status === 200) {
        return response
    }
    console.log('Status NotOk: ', url, '\n', response);
    return false;
}

const _get = async (url) => _callApi(url, {method: 'get'});

const callLoginApi = async (email, password) => _callApi('login', {
    method: 'post',
    body: JSON.stringify({
        Email: email,
        Password: password
    })
})

const callAuthcheckApi = async () => _get('authcheck')
const callUpgradeApi = async () => _get('upgrade')
const callUpgradeCheckApi = async () => _get('upgradecheck')
const callDashboardApi = async () => _get('dashboard')
const callLogoutApi = async () => _get('logout')

const callApi = async () => {
    const response = await fetch('/api');
    if (response.status !== 200) throw Error(response);
    return response;
};


export {
    callDashboardApi,
    callApi,
    callAuthcheckApi,
    callLoginApi,
    RouterHack,
    callUpgradeApi,
    callUpgradeCheckApi,
    callLogoutApi
}