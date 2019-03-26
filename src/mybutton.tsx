import * as React from 'react';


interface IMyButton {
	onHeaderClick?(value);
}

export class MyButton extends React.Component<IMyButton> {
  public handleClick = () => {
      if (this.props.onHeaderClick)
        this.props.onHeaderClick(this.props.children);
  }

  public render() {
    return (
      <button onClick={this.handleClick}>
        {this.props.children}
      </button>
    );
  }
}

export default MyButton;