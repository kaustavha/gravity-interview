import React from 'react';
import { Redirect } from "react-router-dom";

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

    handleSubmit(event) {
        if (this._isInputValid()) {
            this.callApi()
                .then(() => this.setState({loginSuccess: true}))
                .catch((err) => {
                    alert('login failed');
                    this.setState({redirectToReferrer: true});
                });
        } else {
            this.setState({redirectToReferrer: true});
        }
    }

    callApi = async () => {
    const response = await fetch('/api');
    if (response.status !== 200) throw Error(response);
    return body;
    };

    render() {
        let { redirectToReferrer, loginSuccess } = this.state;

        if (redirectToReferrer) {
            return <Redirect to='/login' />;
        } else if (loginSuccess) {
            return <Redirect to='/dashboard' push />;
        }
        return (
            <form className="login-form" onSubmit={this.handleSubmit}>
                <h1>Sign Into Your Account</h1>
        
                <div>
                    <label htmlFor="email">Email Address</label>
                    <input type="email" id="email" className="field" name="email" value={this.state.value} onChange={this.handleChange} />
                </div>
        
                <div>
                    <label htmlFor="password">Password</label>
                    <input type="password" id="password" className="field" name="password" onChange={this.handleChange} />
                </div>
                <input type="button" value="Login to my Dashboard" className="button block" onClick={this.handleSubmit} > 

            </input>
            </form>
        )
    }
}