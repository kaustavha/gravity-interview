import React from 'react';
import {
    callAuthcheckApi,
    callLoginApi
} from "../util/api";
import RouterHack from "./RouterHack";

export default class LoginForm extends React.Component {
    constructor(props) {
        super(props)
        this._isMounted = false
        this.activeReq = false
        this.state = {
            email: 'email',
            password: 'password',
            redirectToReferrer: false,
            loginSuccess: false,
        }
        this.handleChange = this.handleChange.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
    }

    componentWillMount() {
        // Re route already logged in users to dashboard
        this._isMounted = true
        return callAuthcheckApi().then(res => {
            if (res) {
                this.setState({loginSuccess: true})
            }
        })
    }

    componentWillUnmount() {
        // Set unmounted to prevent async leaks on multiple login button presses
        this._isMounted = false
    }

    handleChange(event) {
        const target = event.target
        const name = target.name
        const value = target.value

        this.setState({[name]: value})
    }

    _isInputValid() {
        const email = this.state.email
        const pass = this.state.password
        return (email.length > 0 && pass.length > 0)
    }

    /**
     * handleSubmit: Handles clicking the submit key and logs in user against server
     *  only allow 1 active req, otherwise we can 
     *  smash the login button, queue up a bunch of 
     *  login reqs and then we see a delay in loading 
     *  dashboard state till they all clear
     */
    handleSubmit() {
        if (this._isInputValid() && this._isMounted && !this.activeReq) {
            this.activeReq = true
            return callLoginApi(this.state.email, this.state.password).then(res => {
                if (res && this._isMounted) {
                    this.activeReq = false
                    this.setState({loginSuccess: true})
                }
            })
        }
    }

    render() {
        return (
            <div>
                <RouterHack redirectToReferrer={this.state.redirectToReferrer} loginSuccess={this.state.loginSuccess} />
                <form className="login-form" onSubmit={this.handleSubmit}>
                    <h1>Sign Into Your Account</h1>
            
                    <div>
                        <label htmlFor="email">Email Address</label>
                        <input type="email" autoComplete="username emai" id="email" className="field" name="email" value={this.state.value} onChange={this.handleChange} />
                    </div>
            
                    <div>
                        <label htmlFor="password">Password</label>
                        <input type="password" autoComplete="current-password" id="password" className="field" name="password" onChange={this.handleChange} />
                    </div>
                    <input type="button" value="Login to my Dashboard" className="button block" onClick={this.handleSubmit} />
                </form>
            </div>
        )
    }
}
