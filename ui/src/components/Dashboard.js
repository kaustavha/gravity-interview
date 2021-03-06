import React from 'react';

import {
    callAuthcheckApi,
    callDashboardApi,
    callUpgradeApi,
    callLogoutApi
} from "../util/api";

import RouterHack from "./RouterHack"

export default class Dashboard extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            loginSuccess: false,
            redirectToReferrer: false,
            userCount: 0,
            isAccountUpgraded: false,
            userLimit: 100,
            inProgress: false,
            currentlyUpdatingDB: false
        }

        this.updateDashboard = this.updateDashboard.bind(this)
        this.handleLogout = this.handleLogout.bind(this)
        this.handleUpgrade = this.handleUpgrade.bind(this)
        this.setInitialDashboardState = this.setInitialDashboardState.bind(this)
        this.timer = null;
    }

    componentWillMount() {
        // Reroute non-authorized users to login
        return callAuthcheckApi().then(this.setInitialDashboardState)
    }

    componentWillUnmount() {
        clearTimeout(this.timer)
    }

    setInitialDashboardState(res) {
        if (res) {
            return callDashboardApi().then(resJson => {
                if (resJson) {
                    this.setState({
                        isAccountUpgraded: resJson.IsUpgraded,
                        userLimit: resJson.MaxUsers,
                        userCount: resJson.Users
                    });
                    this.updateDashboard()
                }
            })
        } else {
            this.setState({redirectToReferrer: true})
        }
    }

    updateDashboard() {
        return callDashboardApi().then(resJson => {
            if (resJson) {
                if (resJson.Users < this.state.userLimit) {
                    this.setState({
                        userCount: resJson.Users,
                        currentlyUpdatingDB: true
                    });
                    this.timer = setTimeout(this.updateDashboard.bind(this), 1000)
                } else {
                    this.setState({
                        userCount: resJson.Users,
                        currentlyUpdatingDB: false
                    });
                }
            } else {
                clearTimeout(this.timer)
                this.setState({
                    redirectToReferrer: true
                })
                console.log('update db fail', this, resJson)
            }
        })
    }

    handleLogout() {
        return callLogoutApi().then(() => {
            this.setState({
                redirectToReferrer: true
            })
        })
    }

    handleUpgrade() {
        return callUpgradeApi().then(resJson => {
            if (resJson) {
                this.setState({
                    isAccountUpgraded: resJson.IsUpgraded,
                    userLimit: resJson.MaxUsers,
                    userCount: resJson.Users
                })
                if (!this.state.currentlyUpdatingDB) {
                    this.updateDashboard() // begin refreshing dashboard again
                }
            }
        })
    }

    render() {
        let pct = (this.state.userCount / this.state.userLimit) * 100;

        return (
            <div>
                <RouterHack redirectToReferrer={this.state.redirectToReferrer} loginSuccess={this.state.loginSuccess}/>
                {
                    this.state.userCount >= this.state.userLimit && !this.state.isAccountUpgraded && 
                    <div className="alert is-error">You have exceeded the maximum number of users for your account, please upgrade your plan to increaese the limit.</div>
                }

                {
                    this.state.isAccountUpgraded &&
                    <div className="alert is-success">Your account has been upgraded successfully!</div>
                }
                

                <div className="plan">
                    <header>Startup Plan - $100/Month</header>

                    <div className="plan-content">
                        <div className="progress-bar">
                            <div style={{width: `${pct}%`}} className="progress-bar-usage"></div>
                        </div>

                        <h3>Users: {this.state.userCount}/{this.state.userLimit}</h3>
                    </div>

                    <footer>
                        <button className="button is-error" onClick={this.handleLogout}>Log Out</button>
                        {
                            !this.state.isAccountUpgraded && this.state.userCount >= this.state.userLimit &&
                            <button className="button is-success" onClick={this.handleUpgrade}>Upgrade to Enterprise Plan</button>
                        }
                        
                    </footer>
                </div>
            </div>
        )
    }
}
