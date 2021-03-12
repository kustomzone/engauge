import * as React from "react";
import { Admin, Resource } from 'react-admin';
import { Route } from 'react-router-dom';
import dataProvider from './dataProvider';
import FunctionsIcon from '@material-ui/icons/Functions';
import SettingsIcon from '@material-ui/icons/Settings';
import WebAssetIcon from '@material-ui/icons/WebAsset';
import WebIcon from '@material-ui/icons/Web';
import Category from '@material-ui/icons/Category';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import { createMuiTheme } from '@material-ui/core/styles';
import authProvider from './authProvider';

import { EndpointList, EndpointCreate, EndpointShow } from './endpoint';
import { PropertyList, PropertyCreate, PropertyShow } from './property';
import { OriginList, OriginShow } from './origin';
import { EntityList, EntityShow } from './entity';
import { SummaryList, SummaryShow } from "./summary";
import { SettingsList, SettingsEdit} from "./settings";
// import SettingsEdit from "./settings/index";
// import AppLayout from "./layout";

const theme = createMuiTheme({
     palette: {
          type: 'light',
     },
});

const App = () => (
     <Admin authProvider={authProvider} disableTelemetry theme={theme} dataProvider={dataProvider}
          // customRoutes={[
          //      <Route
          //           key="admin-settings"
          //           path="/admin-settings"
          //           component={SettingsEdit.edit}
          //      />
          // ]}
          // appLayout={AppLayout}
     >
          <Resource name="summaries" list={SummaryList} show={SummaryShow} icon={FunctionsIcon} />
          <Resource name="endpoint" list={EndpointList} create={EndpointCreate} show={EndpointShow} icon={SettingsInputAntennaIcon} />
          <Resource name="origin" list={OriginList} show={OriginShow} icon={WebAssetIcon} />
          <Resource name="entity" options={{ label: 'Entities' }} list={EntityList} show={EntityShow} icon={WebIcon} />
          <Resource name="properties" list={PropertyList} show={PropertyShow} icon={Category} />
          <Resource name="settings" list={SettingsList} edit={SettingsEdit} icon={SettingsIcon} />
     </Admin>
);

export default App;
