import React from 'react';

export default class Dashboard extends React.Component {
    constructor(props) {
        super(props);
    }
    render() {
        return (
            <div>
                <div className="alert is-error">You have exceeded the maximum number of users for your account, please upgrade your plan to increaese the limit.</div>
                <div className="alert is-success">Your account has been upgraded successfully!</div>

                <div className="plan">
                    <header>Startup Plan - $100/Month</header>

                    <div className="plan-content">
                        <div className="progress-bar">
                            <div style={{width: 35 + '%'}} className="progress-bar-usage"></div>
                        </div>

                        <h3>Users: 35/100</h3>
                    </div>

                    <footer>
                        <button className="button is-success">Upgrade to Enterprise Plan</button>
                    </footer>
                </div>
            </div>
        )
    }
}