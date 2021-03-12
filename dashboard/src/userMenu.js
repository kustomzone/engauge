import React, { Component } from 'react';
import { connect } from 'react-redux';
import { crudGetOne, UserMenu, MenuItemLink } from 'react-admin';
import SettingsIcon from '@material-ui/icons/Settings';

class EngaugeUserMenuView extends Component {
    componentDidMount() {
        this.fetchSettings();
    }

    fetchSettings = () => {
        this.props.crudGetOne(
            // resource
            'settings',
            // resource-id
            'admin-settings',
            // base path
            '/admin-settings',
            // refresh
            false
        );
    };

    render() {
        const { crudGetOne, settings, ...props } = this.props;

        return (
            <UserMenu label={settings ? settings.nickname : 'my settings nickname'} {...props}>
                <MenuItemLink
                    to="/admin-settings"
                    primaryText="Settings"
                    leftIcon={<SettingsIcon />}
                />
            </UserMenu>
        );
    }
}

const mapStateToProps = state => {
    const resource = 'settings';
    const id = 'admin-settings';
    return {
        settings: state.admin.resources[resource]
            ? state.admin.resources[resource].data[id]
            : null
    };
};

const EngaugeUserMenu = connect(
    mapStateToProps,
    { crudGetOne }
)(EngaugeUserMenuView);

export default EngaugeUserMenu;