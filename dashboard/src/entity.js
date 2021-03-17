import * as React from "react";
import { List, Datagrid, TextField, ArrayField, NumberField, TabbedShowLayout, Tab, Show, ShowButton } from 'react-admin';

export const EntityList = props => (
    <List {...props}>
        <Datagrid>
            <TextField label="Entity Type" source="entityType" />
            <TextField label="Entity ID" source="entityID" />
            <ShowButton />
        </Datagrid>
    </List>
);

export const EntityShow = (props) => (
    <Show {...props}>
        <TabbedShowLayout>
            <Tab label="Details">
                <TextField label="ID" source="id" />
                <TextField label="Entity Type" source="entityType" />
                <TextField label="Entity ID" source="entityID" />
            </Tab>
            <Tab label="All-Time Stats">
                <ArrayField label="Actions Statistics" source="stats.allTime.stats.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="stats.allTime.stats.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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
            <ArrayField label="Actions Statistics" source="stats.hourly.stats.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="stats.hourly.stats.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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
            <ArrayField label="Actions Statistics" source="stats.daily.stats.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="stats.daily.stats.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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
            <ArrayField label="Actions Statistics" source="stats.weekly.stats.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="stats.weekly.stats.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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
            <ArrayField label="Actions Statistics" source="stats.monthly.stats.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="stats.monthly.stats.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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
            <ArrayField label="Actions Statistics" source="stats.quarterly.stats.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="stats.quarterly.stats.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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
            <ArrayField label="Actions Statistics" source="stats.yearly.stats.actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Statistics" source="stats.yearly.stats.entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entity Type" source="value" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                    </Datagrid>
                </ArrayField>
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