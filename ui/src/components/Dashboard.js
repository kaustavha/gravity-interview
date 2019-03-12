import React from 'react';

import {
    callAuthcheckApi,
    callDashboardApi,
    callUpgradeApi,
    callUpgradeCheckApi,
    callLogoutApi
} from "../util/api";

import RouterHack from "./RouterHack"

export default class Dashboard extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            loginSuccess: false,
            redirectToReferrer: false,
            userCount: 0,
            isAccountUpgraded: false,
            userLimit: 100,
            inProgress: false,
            currentlyUpdatingDB: false
        }
    }

    componentDidMount() {
        // Reroute non-authorized users to login
        callAuthcheckApi().then(res => {
            if (res) {
                callUpgradeCheckApi().then(resJson => {
                    if (resJson) {
                        this.setState({
                            isAccountUpgraded: resJson.IsUpgraded,
                            userLimit: resJson.MaxUsers,
                            userCount: resJson.Users
                        })
                        this.updateDashboard();
                    }
                })
            } else {
                this.setState({redirectToReferrer: true});
            }
        });
    }

    updateDashboard = async () => {
        let resJson = await callDashboardApi()

        if (resJson) {
            if (resJson.userCount < this.state.userLimit) {
                this.setState({
                    userCount: resJson.userCount,
                    currentlyUpdatingDB: true
                });
                setTimeout(this.updateDashboard.bind(this), 1000);
            } else {
                this.setState({
                    userCount: resJson.userCount,
                    currentlyUpdatingDB: false
                });
            }
        } else {
            this.setState({
                redirectToReferrer: true
            })
            console.log('update db fail')
        }
    }

    handleLogout = async () => {
        await callLogoutApi();
        this.setState({
            redirectToReferrer: true
        });
    }

    handleUpgrade = async () => {
        let resJson = await callUpgradeApi()
        if (resJson) {
            this.setState({
                isAccountUpgraded: resJson.IsUpgraded,
                userLimit: resJson.MaxUsers,
                userCount: resJson.Users
            })
            if (!this.state.currentlyUpdatingDB) {
                this.updateDashboard(); // begin refreshing dashboard again
            }
        }
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