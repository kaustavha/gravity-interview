import React, { Component } from 'react';
import {
  BrowserRouter,
  Route,
} from "react-router-dom";

import LoginForm from './LoginForm';
import Dashboard from './Dashboard';
import './App.css';

class App extends Component {
  render() {
    return (
      <BrowserRouter>
          <div>
            <header>
              <meta charSet="UTF-8" />
              <title>User Management Challenge</title>
              <meta name="Kostco" content="Gravitational, Inc." />
            </header>

            <Route exact path="/" component={LoginForm} />
            <Route path="/login" component={LoginForm} />
            <Route path="/dashboard" component={Dashboard} />
          </div>

        </BrowserRouter>
    )
  }
}

export default App;
