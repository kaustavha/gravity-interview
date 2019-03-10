import React from 'react';
import {
    callAuthcheckApi,
    callLoginApi,
    RouterHack
} from "./api";

export default class LoginForm extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            email: 'email',
            password: 'password',
            redirectToReferrer: false,
            loginSuccess: false
        };
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);

        callAuthcheckApi().then(res => {
            if (res) {
                this.setState({loginSuccess: true});
            }
        })
    }

    handleChange(event) {
        const target = event.target;
        const name = target.name;
        const value = target.value;

        this.setState({[name]: value});
    }

    _isInputValid() {
        let email = this.state.email,
            pass = this.state.password;
        return (email.length > 0 && pass.length > 0);
    }

    handleSubmit = async () => {
        if (this._isInputValid()) {
            const res = await callLoginApi(this.state.email, this.state.password);
            if (res) {
                this.setState({loginSuccess: true});
            } else {
                this.setState({redirectToReferrer: true});
            }
        }
    }

    render() {
        return (
            <div>
                <RouterHack redirectToReferrer={this.state.redirectToReferrer} loginSuccess={this.state.loginSuccess}/>
                <form className="login-form" >
                    <h1>Sign Into Your Account</h1>
            
                    <div>
                        <label htmlFor="email">Email Address</label>
                        <input type="email" id="email" className="field" name="email" value={this.state.value} onChange={this.handleChange} />
                    </div>
            
                    <div>
                        <label htmlFor="password">Password</label>
                        <input type="password" id="password" className="field" name="password" onChange={this.handleChange} />
                    </div>
                    <input type="button" value="Login to my Dashboard" className="button block" onClick={this.handleSubmit} />
                </form>
            </div>
        )
    }
}