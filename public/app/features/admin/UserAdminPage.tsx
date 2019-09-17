import React, { PureComponent } from 'react';
import { hot } from 'react-hot-loader';
import { connect } from 'react-redux';
import { NavModel } from '@grafana/data';
import { getNavModel } from 'app/core/selectors/navModel';
import { getRouteParamsId } from 'app/core/selectors/location';
import Page from 'app/core/components/Page/Page';
import { UserProfile } from './UserProfile';
import { StoreState, UserDTO } from 'app/types';
import { loadUserProfile } from './state/actions';

interface Props {
  navModel: NavModel;
  userId: number;
  user: UserDTO;

  loadUserProfile: typeof loadUserProfile;
}

interface State {
  isLoading: boolean;
}

export class UserAdminPage extends PureComponent<Props, State> {
  state = {
    isLoading: true,
  };

  async componentDidMount() {
    const { userId, loadUserProfile } = this.props;
    try {
      await loadUserProfile(userId);
    } finally {
      this.setState({ isLoading: false });
    }
  }

  render() {
    const { navModel, user } = this.props;
    const { isLoading } = this.state;

    return (
      <Page navModel={navModel}>
        <Page.Contents isLoading={isLoading}>
          <UserProfile user={user} />
        </Page.Contents>
      </Page>
    );
  }
}

const mapStateToProps = (state: StoreState) => ({
  userId: getRouteParamsId(state.location),
  navModel: getNavModel(state.navIndex, 'global-users'),
  user: state.userAdmin.user,
  sessions: state.userAdmin.sessions,
});

const mapDispatchToProps = {
  loadUserProfile,
};

export default hot(module)(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(UserAdminPage)
);
