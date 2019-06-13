import React, { PureComponent } from 'react';
import classNames from 'classnames';

import { QueriesTab } from './QueriesTab';
import { VisualizationTab } from './VisualizationTab';
import { AngularComponent } from '@grafana/runtime';
import { PanelResizer } from './PanelResizer';
import { PanelEditHeader } from '../panel_editor/PanelEditHeader';

import { PanelModel } from '../state/PanelModel';
import { DashboardModel } from '../state/DashboardModel';
import { PanelPlugin, PanelPluginMeta } from '@grafana/ui';

interface PanelEditorProps {
  panel: PanelModel;
  dashboard: DashboardModel;
  plugin: PanelPlugin;
  angularPanel?: AngularComponent;
  onTypeChanged: (newType: PanelPluginMeta) => void;
  children: React.ReactNode;
  isEditing: boolean;
  isFullscreen: boolean;
}

export class PanelEditor extends PureComponent<PanelEditorProps> {
  render() {
    const { panel, dashboard, isEditing, isFullscreen, children, angularPanel, onTypeChanged, plugin } = this.props;

    const editorContainerClasses = classNames({
      'panel-editor-container': isEditing,
      'panel-height-helper': !isEditing,
    });

    const editorPanesClasses = classNames({
      'panel-editor-panes': isEditing,
      'panel-height-helper': !isEditing,
    });

    const panelWrapperClass = classNames({
      'panel-wrapper': true,
      'panel-wrapper--edit': isEditing,
      'panel-wrapper--view': isFullscreen && !isEditing,
    });

    return (
      <div className={editorContainerClasses}>
        {panel.isEditing && <PanelEditHeader panel={panel} />}
        <div className={editorPanesClasses}>
          <div className="panel-editor-panes__left">
            <PanelResizer
              isEditing={isEditing}
              panel={panel}
              render={styles => (
                <div className={panelWrapperClass} style={styles}>
                  {children}
                </div>
              )}
            />
            {panel.isEditing && (
              <div className="panel-editor-container__query">
                <QueriesTab panel={panel} dashboard={dashboard} />
              </div>
            )}
          </div>
          {panel.isEditing && (
            <div className="panel-editor-panes__right">
              <VisualizationTab
                panel={panel}
                dashboard={dashboard}
                plugin={plugin}
                onTypeChanged={onTypeChanged}
                angularPanel={angularPanel}
              />
            </div>
          )}
        </div>
      </div>
    );
  }
}
