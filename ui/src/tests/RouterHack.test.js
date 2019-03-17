import React from 'react';
import Dashboard from '../components/Dashboard'
import { mount, configure } from 'enzyme';
import fetch from './__mocks__/fetch';
import sinon from 'sinon';
import Adapter from 'enzyme-adapter-react-16';

global.fetch = fetch;
configure({ adapter: new Adapter() });

it('routes to login page on redirectToReferrer=true', () => {
  sinon.spy(Dashboard.prototype, 'componentWillMount')
  mount(
    <Dashboard/>
  )
  expect(Dashboard.prototype.componentWillMount).toHaveProperty('callCount', 1)
})