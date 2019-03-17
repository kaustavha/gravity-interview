
import React from 'react';
import { Redirect } from "react-router-dom";
// TODO FIXME - prolly a bad way to do routing but it was quick
export default class RouterHack extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            redirectToReferrer: false,
            loginSuccess: false
        }
    }
    shouldComponentUpdate(nextProps) {
        if (nextProps.loginSuccess !== this.state.loginSuccess
            || nextProps.redirectToReferrer !== this.state.redirectToReferrer) {

            this.setState({
                redirectToReferrer: nextProps.redirectToReferrer,
                loginSuccess: nextProps.loginSuccess
            })
            return true
        }
        return false
    }
    render() {
        let { redirectToReferrer, loginSuccess } = this.state;

        if (redirectToReferrer) {
            return <Redirect to='/login'/>
        }
        if (loginSuccess) {
            return <Redirect to='/dashboard' push />
        }
        return <div/>
    }
}
