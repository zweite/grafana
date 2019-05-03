import React, { PureComponent } from 'react';
import uniqueId from 'lodash/uniqueId';

export interface Props {
  label: string;
  checked: boolean;
  className?: string;
  labelClass?: string;
  switchClass?: string;
  transparent?: boolean;
  description?: string;
  onChange: (event?: React.SyntheticEvent<HTMLInputElement>) => void;
}

export interface State {
  id: any;
}

export class Switch extends PureComponent<Props, State> {
  state = {
    id: uniqueId(),
  };

  internalOnChange = (event: React.FormEvent<HTMLInputElement>) => {
    event.stopPropagation();

    this.props.onChange(event);
  };

  render() {
    const { switchClass = '', label, checked, transparent, className, description } = this.props;

    const labelId = this.state.id;
    const labelClassName = `gf-form-label ${transparent ? 'gf-form-label--transparent' : ''} pointer`;
    const switchClassName = `gf-form-switch ${switchClass} ${transparent ? 'gf-form-switch--transparent' : ''}`;

    return (
      <div className="gf-form-switch-container-react form-field">
        <label htmlFor={labelId} className={`gf-form-switch-container ${className || ''}`}>
          {label && <div className={labelClassName}>{label}</div>}
          {description && <div className="form-field-desc">{description}</div>}
          <div className={switchClassName}>
            <input id={labelId} type="checkbox" checked={checked} onChange={this.internalOnChange} />
            <span className="gf-form-switch__slider" />
          </div>
        </label>
      </div>
    );
  }
}
