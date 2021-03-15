import * as React from "react";
import { List, Datagrid, TextField, EditButton, Edit, Create, SimpleForm, useNotify, useRefresh, useRedirect, TextInput, SelectInput, ArrayInput, SimpleFormIterator, ReferenceInput, SimpleShowLayout, Show, ArrayField, ReferenceField, ShowButton, BooleanField, BooleanInput, NumberInput } from 'react-admin';

export const SettingsList = props => (
    <List title="Settings" {...props}>
        <Datagrid>
            <EditButton />
        </Datagrid>
    </List>
);

export const SettingsEdit = props => (
    <Edit title="Settings" {...props}>
        <SimpleForm>
            <BooleanInput label="Interactions Storage" source="interactions" />
            <BooleanInput label="Hourly Stats" source="statsToggles.hourly.on" />
            <BooleanInput label="Daily Stats" source="statsToggles.daily.on" />
            <BooleanInput label="Weekly Stats" source="statsToggles.weekly.on" />
            <BooleanInput label="Monthly Stats" source="statsToggles.monthly.on" />
        </SimpleForm>
    </Edit>
);