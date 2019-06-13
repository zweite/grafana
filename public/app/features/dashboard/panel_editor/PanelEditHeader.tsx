// Libraries
import React, { PureComponent } from 'react';

// import { DashNavButton } from '../components/DashNav/DashNavButton';

// Types
import { PanelModel } from '../state';

interface Props {
  panel: PanelModel;
}

export class PanelEditHeader extends PureComponent<Props> {
  onClose = () => {};

  render() {
    const { panel } = this.props;

    return (
      <div className="panel-edit-header">
        <div className="navbar">
          <div className="navbar-edit">
            <button className="navbar-edit__back-btn" onClick={this.onClose}>
              <i className="fa fa-arrow-left" />
            </button>
          </div>

          <div className="navbar-page-btn">{panel.title}</div>

          <div className="navbar__spacer" />
          <div className="navbar-buttons">
            <button className="btn navbar-button navbar-button--settings">
              Graph <i className="fa fa-caret-down" />
            </button>
            <button className="btn navbar-button navbar-button--settings">
              <i className="gicon gicon-cog" />
            </button>
          </div>
        </div>
      </div>
    );
  }
}

export class PanelEditOptionsPane extends PureComponent<Props> {
  render() {
    return <div className="panel-editor-options-pane">Settings</div>;
  }
}
