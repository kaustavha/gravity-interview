import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import {
  RouterHack
} from "./api";

it('renders without crashing', () => {
  const div1 = document.createElement('div');
  ReactDOM.render(<RouterHack />, div1);
  ReactDOM.unmountComponentAtNode(div1);

  const div = document.createElement('div');
  ReactDOM.render(<App />, div);
  ReactDOM.unmountComponentAtNode(div);
});
