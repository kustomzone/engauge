import * as React from "react";
import { List, Datagrid, ShowButton, TextField, NumberField, TabbedShowLayout, Tab, ArrayField, Create, SimpleForm, useNotify, useRefresh, useRedirect, TextInput, SelectInput, Show } from 'react-admin';

export const PropertyList = props => (
    <List {...props}>
        <Datagrid>
            <TextField label="Name" source="id" />
            <TextField source="type" />
            <ShowButton />
        </Datagrid>
    </List>
);

export const PropertyCreate = props => {
    const notify = useNotify();
    const refresh = useRefresh();
    const redirect = useRedirect();

    const onSuccess = ({ data }) => {
        notify(`Property ${data.name} of type ${data.type} has been created.`)
        redirect('/properties');
        refresh();
    };

    const onFailure = ({ error }) => {
        notify(`ERROR: ${error}`)
        redirect('/properties');
        refresh();
    };

    return (
        <Create {...props} onSuccess={onSuccess} onFailure={onFailure}>
            <SimpleForm>
                <TextInput source="name" />
                <SelectInput source="type" choices={[
                    { id: 'string', name: 'string' },
                    { id: 'number', name: 'number' },
                    { id: 'number-array', name: 'number-array' },
                    { id: 'string-array', name: 'string-array' },
                ]} />
            </SimpleForm>
        </Create>
    );
};

export const PropertyShow = props => (
    <Show {...props}>
        <TabbedShowLayout>
            <Tab label="Details">
                <TextField label="Name" source="id" />
                <TextField label="Type" source="type" />
            </Tab>
            <Tab label="All-Time Stats">
                <TextField label="Average Value" source="stats.mean" />
                <TextField label="Most Occurring Value" source="stats.mode" />
                <ArrayField label="Value Statistics" source="stats.values">
                    <Datagrid>
                        <TextField label="Value" source="value" />
                        <NumberField label="Total" source="count" />
                        <NumberField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
            </Tab>
            <Tab label="Hourly Stats">
                <TextField label="Average Value" source="hourlyStats.stats.mean" />
                <TextField label="Most Occurring Value" source="hourlyStats.stats.mode" />
                <ArrayField label="Value Statistics" source="hourlyStats.stats.values">
                    <Datagrid>
                        <TextField label="Value" source="value" />
                        <NumberField label="Total" source="count" />
                        <NumberField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
            </Tab>
            <Tab label="Daily Stats">
                <TextField label="Average Value" source="dailyStats.stats.mean" />
                <TextField label="Most Occurring Value" source="dailyStats.stats.mode" />
                <ArrayField label="Value Statistics" source="dailyStats.stats.values">
                    <Datagrid>
                        <TextField label="Value" source="value" />
                        <NumberField label="Total" source="count" />
                        <NumberField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
            </Tab>
            <Tab label="Weekly Stats">
                <TextField label="Average Value" source="weeklyStats.stats.mean" />
                <TextField label="Most Occurring Value" source="weeklyStats.stats.mode" />
                <ArrayField label="Value Statistics" source="weeklyStats.stats.values">
                    <Datagrid>
                        <TextField label="Value" source="value" />
                        <NumberField label="Total" source="count" />
                        <NumberField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
            </Tab>
            <Tab label="Monthly Stats">
                <TextField label="Average Value" source="monthlyStats.stats.mean" />
                <TextField label="Most Occurring Value" source="monthlyStats.stats.mode" />
                <ArrayField label="Value Statistics" source="monthlyStats.stats.values">
                    <Datagrid>
                        <TextField label="Value" source="value" />
                        <NumberField label="Total" source="count" />
                        <NumberField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
            </Tab>
            <Tab label="Quarterly Stats">
                <TextField label="Average Value" source="quarterlyStats.stats.mean" />
                <TextField label="Most Occurring Value" source="quarterlyStats.stats.mode" />
                <ArrayField label="Value Statistics" source="quarterlyStats.stats.values">
                    <Datagrid>
                        <TextField label="Value" source="value" />
                        <NumberField label="Total" source="count" />
                        <NumberField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
            </Tab>
            <Tab label="Yearly Stats">
                <TextField label="Average Value" source="yearlyStats.stats.mean" />
                <TextField label="Most Occurring Value" source="yearlyStats.stats.mode" />
                <ArrayField label="Value Statistics" source="yearlyStats.stats.values">
                    <Datagrid>
                        <TextField label="Value" source="value" />
                        <NumberField label="Total" source="count" />
                        <NumberField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
            </Tab>
        </TabbedShowLayout>
    </Show>
);