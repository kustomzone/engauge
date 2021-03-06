import * as React from "react";
import { List, Datagrid, ShowButton, TextField, NumberField, TabbedShowLayout, Tab, ArrayField, Create, SimpleForm, useNotify, useRefresh, useRedirect, TextInput, SelectInput, ReferenceInput, Filter, Show, SimpleShowLayout, DateField } from 'react-admin';

export const SummaryList = props => (
    <List {...props}>
        <Datagrid>
            <TextField label="Interval" source="id" />
            <DateField label="Starting Time" source="start" />
            <DateField label="Ending Time" source="end" />
            <ShowButton />
        </Datagrid>
    </List>
);

export const SummaryShow = props => (
    <Show {...props}>
        <TabbedShowLayout>
            <Tab label="Details">
                <TextField label="Interval" source="id" />
                <DateField label="Starting Time" source="start" />
                <DateField label="Ending Time" source="end" />
                <NumberField label="Total Interactions" source="total" />
            </Tab>
            <Tab label="General Stats">
                <ArrayField label="Action Stats" source="actionStats.values">
                    <Datagrid>
                        <TextField label="Action" source="value" />
                        <TextField label="Total" source="count" />
                        <TextField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Origin Type Stats" source="originTypeStats.values">
                    <Datagrid>
                        <TextField label="Origin Type" source="value" />
                        <TextField label="Total" source="count" />
                        <TextField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Entity Type Stats" source="entityTypeStats.values">
                    <Datagrid>
                        <TextField label="Entiity Type" source="value" />
                        <TextField label="Total" source="count" />
                        <TextField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="User Type Stats" source="userTypeStats.values">
                    <Datagrid>
                        <TextField label="User Type" source="value" />
                        <TextField label="Total" source="count" />
                        <TextField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Device Type Stats" source="deviceTypeStats.values">
                    <Datagrid>
                        <TextField label="Device Type" source="value" />
                        <TextField label="Total" source="count" />
                        <TextField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
                <ArrayField label="Session Type Stats" source="sessionTypeStats.values">
                    <Datagrid>
                        <TextField label="Session Type" source="value" />
                        <TextField label="Total" source="count" />
                        <TextField label="Percentage" source="percentage" />
                    </Datagrid>
                </ArrayField>
            </Tab>
            <Tab label="Session Stats">
                <ArrayField label="Session Statistics" source="sessionStats">
                    <Datagrid>
                        <TextField label="User Type" source="userType" />
                        <TextField label="Device Type" source="deviceType" />
                        <TextField label="Session Type" source="sessionType" />
                        <NumberField label="Total Interactions" source="count" />
                        <NumberField label="Percentage of Interactions" source="percentage" />
                        <NumberField label="Total Conversions" source="conversions" />
                        <NumberField label="Conversion Rate" source="conversionRate" />
                        <NumberField label="Bounced Sessions" source="bouncedSessions" />
                        <NumberField label="Bounce Rate" source="bounceRate" />
                        <NumberField label="Average Duration (minutes)" source="durationStats.mean" />
                        <NumberField label="Average Interactions Per Session" source="interactionStats.mean" />
                    </Datagrid>
                </ArrayField>
            </Tab>
            <Tab label="Conversion Stats">
                <ArrayField label="Conversion Endpoint Statistics" source="conversionStats">
                    <Datagrid>
                        <TextField label="Conversion Endpoint" source="endpoint" />
                        <NumberField label="Total Revenue" source="value" />
                        <NumberField label="Average Revenue Per Conversion" source="amountStats.mean" />
                    </Datagrid>
                </ArrayField>
            </Tab>
            <Tab label="Unit Metrics">
                <NumberField label="Total Conversions" source="unitMetrics.totalConversions" />
                <NumberField label="Total Revenue" source="unitMetrics.totalRevenue" />
                <NumberField label="ARPU" source="unitMetrics.averageRevenuePerUser" />
                <NumberField label="Average Conversion Amount" source="unitMetrics.amountStats.mean" />
                <NumberField label="Most Occuring Conversion Amount" source="unitMetrics.amountStats.mode" />
            </Tab>
        </TabbedShowLayout>
    </Show>
);