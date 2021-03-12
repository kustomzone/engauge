import * as React from "react";
import { List, Datagrid, TextField, Create, useRedirect, ArrayField, NumberField, TabbedShowLayout, Tab, useRefresh, useNotify, SimpleForm, TextInput, SelectInput, ReferenceInput, Filter, Show, ShowButton } from 'react-admin';

export const OriginList = props => (
    <List {...props}>
        <Datagrid>
            <TextField label="Origin Type" source="originType" />
            <TextField label="Origin ID" source="originID" />
            <ShowButton />
        </Datagrid>
    </List>
);

export const OriginShow = (props) => (
    <Show {...props}>
        <TabbedShowLayout>
            <Tab label="Details">
                <TextField label="ID" source="id" />
                <TextField label="Origin Type" source="originType" />
                <TextField label="Origin ID" source="originID" />
            </Tab>
            <Tab label="All-Time Stats">
                <ArrayField label="Actions Statistics" source="allTimeStats.profile.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="allTimeStats.profile.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="User Type Statistics" source="allTimeStats.profile.userTypeStats.values">
                    <Datagrid>
                        <TextField label="User Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Session Type Statistics" source="allTimeStats.profile.sessionTypeStats.values">
                    <Datagrid>
                        <TextField label="Session Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Device Type Statistics" source="allTimeStats.profile.deviceTypeStats.values">
                    <Datagrid>
                        <TextField label="Device Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Property Statistics" source="allTimeStats.profile.propertyStats.List">
                    <Datagrid>
                        <TextField label="Property" source="name" />
                        <TextField label="Total Interactions" source="stats.total" />
                        <TextField label="Average Value" source="stats.mean" />
                        <TextField label="Most Occuring Value" source="stats.mode" />
                    </Datagrid>
                </ArrayField>
            </Tab>
            <Tab label="Hourly Stats">
            <ArrayField label="Actions Statistics" source="hourlyStats.profile.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="hourlyStats.profile.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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
            <ArrayField label="Actions Statistics" source="dailyStats.profile.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="dailyStats.profile.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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
            <ArrayField label="Actions Statistics" source="weeklyStats.profile.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="weeklyStats.profile.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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
            <ArrayField label="Actions Statistics" source="monthlyStats.profile.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="monthlyStats.profile.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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