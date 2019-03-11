
import React from 'react';
import { Redirect } from "react-router-dom";
// TODO FIXME - prolly a bad way to do routing but it was quick
export default class RouterHack extends React.Component {
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
