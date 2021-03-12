import React from 'react';
import { Layout } from 'react-admin';
import EngaugeAppBar from './appBar';

const AppLayout = props => <Layout {...props} appBar={EngaugeAppBar} />;

export default AppLayout;