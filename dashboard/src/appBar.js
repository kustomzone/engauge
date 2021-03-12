import React from 'react';
import { AppBar } from 'react-admin';
import EngaugeUserMenu from './userMenu';

const EngaugeAppBar = props => <AppBar {...props} userMenu={<EngaugeUserMenu />} />;
export default EngaugeAppBar;