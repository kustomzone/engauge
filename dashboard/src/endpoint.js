import * as React from "react";
import { List, Datagrid, TextField, useRedirect, useRefresh, useNotify, Create, ArrayField, SimpleForm, TextInput, ReferenceInput, Filter, Show, Tab, TabbedShowLayout, ShowButton, NumberField } from 'react-admin';

export const EndpointList = props => (
    <List filters={<EndpointFilter />} {...props}>
        <Datagrid>
            <TextField source="action" />
            <TextField label="Entity Type" source="entityType" />
            <TextField label="Entity ID" source="entityID" />
            <TextField label="Origin Type" source="originType" />
            <TextField label="Origin ID" source="originID" />
            <ShowButton />
        </Datagrid>
    </List>
);

export const EndpointCreate = props => {
    const notify = useNotify();
    const refresh = useRefresh();
    const redirect = useRedirect();

    const onSuccess = ({ data }) => {
        notify(`Endpoint ${data.id} has been created.`)
        redirect('/endpoint');
        refresh();
    };

    const onFailure = ({ error }) => {
        notify(`ERROR: ${error}`)
        redirect('/endpoint');
        refresh();
    };

    return (
        <Create {...props} onSuccess={onSuccess} onFailure={onFailure}>
            <SimpleForm>
                <TextInput required source="action" />
                <TextInput label="Entity Type" source="entityType" />
                <TextInput label="Entity ID" source="entityID" />
                <TextInput label="Origin Type" source="originType" />
                <TextInput label="Origin ID" source="originID" />
            </SimpleForm>
        </Create>
    );
}

const EndpointFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Search" source="q" alwaysOn />
        <ReferenceInput label="Status" source="status" allowEmpty>
            <TextInput source="action" />
            <TextInput source="entityType" />
            <TextInput source="originType" />
        </ReferenceInput>
    </Filter>
);

export const EndpointShow = (props) => {
    return (
        <Show {...props}>
            <TabbedShowLayout>
                <Tab label="Details">
                    <TextField label="ID" source="id" />
                    <TextField source="action" />
                    <TextField label="Entity Type" source="entityType" />
                    <TextField label="Entity ID" source="entityID" />
                    <TextField label="Origin Type" source="originType" />
                    <TextField label="Origin ID" source="originID" />
                </Tab>
                <Tab label="All-Time Stats">
                    <NumberField label="Total Interactions" source="profile.total" />
                    <ArrayField label="User Type Statistics" source="profile.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="profile.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="profile.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="profile.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
                <Tab label="Hourly Stats">
                    <NumberField label="Total Interactions" source="hourlyStats.profile.total" />
                    <ArrayField label="User Type Statistics" source="hourlyStats.profile.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="hourlyStats.profile.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="hourlyStats.profile.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="hourlyStats.profile.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
                <Tab label="Daily Stats">
                    <NumberField label="Total Interactions" source="dailyStats.profile.total" />
                    <ArrayField label="User Type Statistics" source="dailyStats.profile.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="dailyStats.profile.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="dailyStats.profile.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="dailyStats.profile.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
                <Tab label="Weekly Stats">
                    <NumberField label="Total Interactions" source="weeklyStats.profile.total" />
                    <ArrayField label="User Type Statistics" source="weeklyStats.profile.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="weeklyStats.profile.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="weeklyStats.profile.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="weeklyStats.profile.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
                <Tab label="Monthly Stats">
                    <NumberField label="Total Interactions" source="monthlyStats.profile.total" />
                    <ArrayField label="User Type Statistics" source="monthlyStats.profile.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="monthlyStats.profile.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="monthlyStats.profile.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="monthlyStats.profile.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
            </TabbedShowLayout>
        </Show>
    );
};