import React, { FunctionComponent, ReactNode } from 'react';
import classNames from 'classnames';

interface Props {
  children: ReactNode;
  className?: string;
  htmlFor?: string;
  tooltip?: string | React.ReactNode;
  width?: number;
}

export const FormLabel: FunctionComponent<Props> = ({ children, className, htmlFor, tooltip, width, ...rest }) => {
  const classes = classNames('form-label', className);

  return (
    <>
      <label className={classes} {...rest} htmlFor={htmlFor}>
        {children}
      </label>
      {tooltip && <div className="form-label-description">{tooltip}</div>}
    </>
  );
};
