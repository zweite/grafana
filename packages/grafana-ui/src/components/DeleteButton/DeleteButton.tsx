import React, { PureComponent, SyntheticEvent } from 'react';
import styled from 'styled-components';

const Container = styled.span`
  width: 24px;
  direction: rtl;
  display: flex;
  align-items: center;
`;

export  const StyledDeleteButton = styled.a`
  position: absolute;
  transition: opacity 0.1s ease;
  z-index: 0;
`;

const ConfirmDeleteContainer = styled.span<ConfirmDeleteWrapperProps>`
  overflow: hidden;
  width: 145px;
  position: absolute;
  z-index: 1;
  visibility: ${props => (props.showConfirm ? 'visible' : 'hidden')};
`;

interface ConfirmDeleteWrapperProps { showConfirm: boolean; }

export const ConfirmDeleteWrapper = styled.span<ConfirmDeleteWrapperProps>`
  display: flex;
  align-items: flex-start;
  transition: opacity 0.12s ease-in, transform 0.14s ease-in;
  transform: ${props => (props.showConfirm ? 'translateX(0)' : 'translateX(100px)')};
  opacity: ${props => (props.showConfirm ? 1 : 0)};
`;

interface Props {
  onConfirm(): void;
}

interface State {
  showConfirm: boolean;
}

export class DeleteButton extends PureComponent<Props, State> {
  state: State = {
    showConfirm: false,
  };

  onClickDelete = (event: SyntheticEvent) => {
    if (event) {
      event.preventDefault();
    }

    this.setState({
      showConfirm: true,
    });
  };

  onClickCancel = (event: SyntheticEvent) => {
    if (event) {
      event.preventDefault();
    }
    this.setState({
      showConfirm: false,
    });
  };

  render() {
    const { onConfirm } = this.props;
    let showConfirm;
    let showDeleteButton;

    if (this.state.showConfirm) {
      showConfirm = true;
      showDeleteButton = false;
    } else {
      showConfirm = false;
      showDeleteButton = true;
    }

    return (
      <Container>
        {showDeleteButton && (
          <StyledDeleteButton className="btn btn-danger btn-small" onClick={this.onClickDelete}>
            <i className="fa fa-remove" />
          </StyledDeleteButton>
        )}

        <ConfirmDeleteContainer showConfirm={showConfirm}>
          <ConfirmDeleteWrapper showConfirm={showConfirm}>
            <a className="btn btn-small" onClick={this.onClickCancel}>
              Cancel
            </a>
            <a className="btn btn-danger btn-small" onClick={onConfirm}>
              Confirm Delete
            </a>
          </ConfirmDeleteWrapper>
        </ConfirmDeleteContainer>
      </Container>
    );
  }
}
