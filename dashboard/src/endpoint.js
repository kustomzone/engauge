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
                    <NumberField label="Total Interactions" source="stats.allTime.stats.total" />
                    <ArrayField label="User Type Statistics" source="stats.allTime.stats.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="stats.allTime.stats.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="stats.allTime.stats.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="stats.allTime.stats.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
                <Tab label="Hourly Stats">
                    <NumberField label="Total Interactions" source="stats.hourly.stats.total" />
                    <ArrayField label="User Type Statistics" source="stats.hourly.stats.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="stats.hourly.stats.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="stats.hourly.stats.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="stats.hourly.stats.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
                <Tab label="Daily Stats">
                    <NumberField label="Total Interactions" source="stats.daily.stats.total" />
                    <ArrayField label="User Type Statistics" source="stats.daily.stats.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="stats.daily.stats.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="stats.daily.stats.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="stats.daily.stats.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
                <Tab label="Weekly Stats">
                    <NumberField label="Total Interactions" source="stats.weekly.stats.total" />
                    <ArrayField label="User Type Statistics" source="stats.weekly.stats.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="stats.weekly.stats.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="stats.weekly.stats.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="stats.weekly.stats.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
                <Tab label="Monthly Stats">
                    <NumberField label="Total Interactions" source="stats.monthly.stats.total" />
                    <ArrayField label="User Type Statistics" source="stats.monthly.stats.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="stats.monthly.stats.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="stats.monthly.stats.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="stats.monthly.stats.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
                <Tab label="Quarterly Stats">
                    <NumberField label="Total Interactions" source="stats.quarterly.stats.total" />
                    <ArrayField label="User Type Statistics" source="stats.quarterly.stats.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="stats.quarterly.stats.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="stats.quarterly.stats.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="stats.quarterly.stats.propertyStats.List">
                        <Datagrid>
                            <TextField label="Property" source="name" />
                            <TextField label="Total Interactions" source="stats.total" />
                            <TextField label="Average Value" source="stats.mean" />
                            <TextField label="Most Occuring Value" source="stats.mode" />
                        </Datagrid>
                    </ArrayField>
                </Tab>
                <Tab label="Yearly Stats">
                    <NumberField label="Total Interactions" source="stats.yearly.stats.total" />
                    <ArrayField label="User Type Statistics" source="stats.yearly.stats.userTypeStats.values">
                        <Datagrid>
                            <TextField label="User Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Session Type Statistics" source="stats.yearly.stats.sessionTypeStats.values">
                        <Datagrid>
                            <TextField label="Session Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Device Type Statistics" source="stats.yearly.stats.deviceTypeStats.values">
                        <Datagrid>
                            <TextField label="Device Type" source="value" />
                            <NumberField label="Total Interactions" source="count" />
                            <NumberField label="Percentage of Interactions" source="percentage" />
                        </Datagrid>
                    </ArrayField>
                    <ArrayField label="Property Statistics" source="stats.yearly.stats.propertyStats.List">
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