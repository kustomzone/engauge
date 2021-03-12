import React from 'react';
import { Edit, TextInput, BooleanInput, SimpleForm, required } from 'react-admin';

const SettingsEdit = ({ staticContext, ...props }) => {
    return (
        <Edit
            id="admin-settings"
            resource="settings"
            basePath="/admin-settings"
            redirect={false}
            title="Settings"
            {...props}
        >
            <SimpleForm>
                <BooleanInput source="interactions" validate={required()} />
            </SimpleForm>
        </Edit>
    );
};

export default SettingsEdit;