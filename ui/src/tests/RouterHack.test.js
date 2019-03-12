import React from 'react';
import ReactDOM from 'react-dom';
import RouterHack from '../components/RouterHack'

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<RouterHack />, div);
  ReactDOM.unmountComponentAtNode(div);
});
