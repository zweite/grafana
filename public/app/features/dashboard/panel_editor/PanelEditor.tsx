import React, { PureComponent } from 'react';
import classNames from 'classnames';

import { QueriesTab } from './QueriesTab';
import VisualizationTab from './VisualizationTab';
import { GeneralTab } from './GeneralTab';
import { AlertTab } from '../../alerting/AlertTab';

import config from 'app/core/config';
import { store } from 'app/store/store';
import { updateLocation } from 'app/core/actions';
import { AngularComponent } from '@grafana/runtime';

import { PanelModel } from '../state/PanelModel';
import { DashboardModel } from '../state/DashboardModel';
import { Tooltip, PanelPlugin, PanelPluginMeta } from '@grafana/ui';

interface PanelEditorProps {
  panel: PanelModel;
  dashboard: DashboardModel;
  plugin: PanelPlugin;
  angularPanel?: AngularComponent;
  onTypeChanged: (newType: PanelPluginMeta) => void;
}

interface PanelEditorTab {
  id: string;
  text: string;
}

enum PanelEditorTabIds {
  Queries = 'queries',
  Visualization = 'visualization',
  Advanced = 'advanced',
  Alert = 'alert',
}

interface PanelEditorTab {
  id: string;
  text: string;
}

const panelEditorTabTexts = {
  [PanelEditorTabIds.Queries]: 'Queries',
  [PanelEditorTabIds.Visualization]: 'Visualization',
  [PanelEditorTabIds.Advanced]: 'General',
  [PanelEditorTabIds.Alert]: 'Alert',
};

const getPanelEditorTab = (tabId: PanelEditorTabIds): PanelEditorTab => {
  return {
    id: tabId,
    text: panelEditorTabTexts[tabId],
  };
};

export class PanelEditor extends PureComponent<PanelEditorProps> {
  constructor(props: PanelEditorProps) {
    super(props);
  }

  onChangeTab = (tab: PanelEditorTab) => {
    store.dispatch(
      updateLocation({
        query: { tab: tab.id, openVizPicker: null },
        partial: true,
      })
    );
    this.forceUpdate();
  };

  renderCurrentTab(activeTab: string) {
    const { panel, dashboard, onTypeChanged, plugin, angularPanel } = this.props;

    switch (activeTab) {
      case 'advanced':
        return <GeneralTab panel={panel} />;
      case 'queries':
        return <QueriesTab panel={panel} dashboard={dashboard} />;
      case 'alert':
        return <AlertTab angularPanel={angularPanel} dashboard={dashboard} panel={panel} />;
      case 'visualization':
        return (
          <VisualizationTab
            panel={panel}
            dashboard={dashboard}
            plugin={plugin}
            onTypeChanged={onTypeChanged}
            angularPanel={angularPanel}
          />
        );
      default:
        return null;
    }
  }

  render() {
    const { plugin } = this.props;
    let activeTab: PanelEditorTabIds = store.getState().location.query.tab || PanelEditorTabIds.Queries;

    const tabs: PanelEditorTab[] = [
      getPanelEditorTab(PanelEditorTabIds.Queries),
      getPanelEditorTab(PanelEditorTabIds.Visualization),
      getPanelEditorTab(PanelEditorTabIds.Advanced),
    ];

    // handle panels that do not have queries tab
    if (plugin.meta.skipDataQuery) {
      // remove queries tab
      tabs.shift();
      // switch tab
      if (activeTab === PanelEditorTabIds.Queries) {
        activeTab = PanelEditorTabIds.Visualization;
      }
    }

    if (config.alertingEnabled && plugin.meta.id === 'graph') {
      tabs.push(getPanelEditorTab(PanelEditorTabIds.Alert));
    }

    return <div className="panel-editor-container__query">{this.renderCurrentTab(activeTab)}</div>;
  }
}
