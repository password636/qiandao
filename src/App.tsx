import * as React from 'react';
import './App.css';
import {RemoteClient} from "./httprequest"
import "./style.css"
import {MyButton} from "./mybutton"

class App extends React.Component {

  public state = {
    signed: [] as string[],
    unsigned: [] as JSX.Element[]
	}

	constructor(props) {
		super(props)
	}

	public componentDidMount() {
		this.view();
  }

  
      
	private  async view() {
		const client = new RemoteClient({url:"/api/"})
		const names = await client.get()
    const sendTime = async (e) => {
      if(window.confirm(e+'，你确定要签到吗？')){
        const nodeDate = require('date-and-time');
        const now = nodeDate.format(new Date(), 'YYYY/MM/DD hh:mm:ss');
        const newnames = await client.sign({name:e, date:now})

        const rowUnsigned:JSX.Element[] = []
        if (newnames.unsigned) {
          for (const i of newnames.unsigned) {
            rowUnsigned.push(<MyButton onHeaderClick={sendTime}>{i}</MyButton>)
          }
        } else {
          rowUnsigned.push(<input type="button" value="已经全部签到"/>)
        }

        const rowSigned:JSX.Element[] = []
        for (const i of newnames.signed) {
          console.log("fff", i)
          rowSigned.push(<input type="button" value={i}/>)
          rowSigned.push(<br/>)
        }

        this.setState({signed : <div>{rowSigned}</div>, unsigned : <div>{rowUnsigned}</div>});
     }
    }


    const rows:JSX.Element[] = []
      for (const i of names.unsigned) {
        console.log("ddd", i)
        rows.push(<MyButton onHeaderClick={sendTime}>{i}</MyButton>)
    }
    
		this.setState({signed : names.signed, unsigned : <div>{rows}</div>});
		
	}

  public render() {

    return (
    <div>

      <h1>
        未签到名单
      </h1>
      
      {this.state.unsigned}
      

      <h1>
        已签到名单
      </h1>
      {this.state.signed}


    </div>
    );
  }
}

export default App;
