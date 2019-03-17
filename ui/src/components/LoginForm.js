import React from 'react';
import {
    callAuthcheckApi,
    callLoginApi
} from "../util/api";
import RouterHack from "./RouterHack";

export default class LoginForm extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            email: 'email',
            password: 'password',
            redirectToReferrer: false,
            loginSuccess: false
        }
        this.handleChange = this.handleChange.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
        // this._isInputValid = this._isInputValid.bind(this)
    }

    componentDidMount() {
        // Re route already logged in users to dashboard
        return callAuthcheckApi().then(res => {
            if (res) {
                this.setState({loginSuccess: true})
            }
        })
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

    handleSubmit() {
        if (this._isInputValid()) {
            return callLoginApi(this.state.email, this.state.password).then(res => {
                if (res) {
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